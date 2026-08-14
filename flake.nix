{
  description = "cool.sacha.house";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs =
    { nixpkgs, self }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          app = pkgs.buildGoModule {
            pname = "cool-sacha-house";
            version = "0.0.0";
            src = self;
            vendorHash = "sha256-PZlBt0TmxQEkpxoVTV3Ob9OqQ5c0BeP8m2rnhmu1cec=";
            subPackages = [ "cmd" ];
            env.CGO_ENABLED = "0";
            ldflags = [
              "-s"
              "-w"
            ];
            postInstall = ''
              mv "$out/bin/cmd" "$out/bin/cool-sacha-house"
            '';
          };
          image = pkgs.dockerTools.buildLayeredImage {
            name = "cool-sacha-house";
            tag = "latest";
            created = "1970-01-01T00:00:01Z";
            extraCommands = ''
              mkdir -p app var/db
              cp ${app}/bin/cool-sacha-house app/
              cp -r ${./assets} app/assets
              cp -r ${./views} app/views
            '';
            config = {
              Entrypoint = [ "/app/cool-sacha-house" ];
              WorkingDir = "/app";
              Env = [
                "PORT=7883"
                "DB_URL=/var/db/prod.db"
                "ENCRYPTION_KEY=ABCDEFGHIJKLMNOPQRSTUVWXYZ"
                "VERSION=v0.0.0+nix"
                "COMMIT_SHA=nix"
              ];
              ExposedPorts."7883/tcp" = { };
              Volumes."/var/db" = { };
            };
          };
        in
        {
          default = app;
          inherit app;
          dockerImage = image;
        }
      );

      checks = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          app = self.packages.${system}.app;
          goCheck =
            name: command:
            app.overrideAttrs (_: {
              pname = name;
              doCheck = false;
              buildPhase = command;
              installPhase = "touch $out";
              postInstall = null;
            });
        in
        {
          inherit app;
          dockerImage = self.packages.${system}.dockerImage;

          actionlint =
            pkgs.runCommand "actionlint"
              {
                nativeBuildInputs = [ pkgs.actionlint ];
                src = self;
              }
              ''
                cd "$src"
                actionlint -config-file .github/actionlint.yaml .github/workflows/*.yml
                touch "$out"
              '';

          format =
            pkgs.runCommand "format"
              {
                nativeBuildInputs = [
                  pkgs.go
                  pkgs.nixfmt
                ];
                src = self;
              }
              ''
                cd "$src"
                test -z "$(gofmt -l .)"
                nixfmt --check flake.nix
                touch "$out"
              '';

          go-vet = goCheck "go-vet" "go vet ./...";
          tests = goCheck "go-tests" "go test ./...";

          http =
            pkgs.runCommand "image-http-test"
              {
                nativeBuildInputs = [
                  pkgs.curl
                  pkgs.jq
                ];
              }
              ''
                mkdir archive root db
                tar -xf ${self.packages.${system}.dockerImage} -C archive
                jq -r '.[0].Layers[]' archive/manifest.json | while read -r layer; do
                  tar -xf "archive/$layer" -C root
                done
                cd root/app
                PORT=18783 DB_URL="$NIX_BUILD_TOP/db/test.db" \
                  ./cool-sacha-house >server.log 2>&1 &
                server_pid=$!
                trap 'kill "$server_pid" 2>/dev/null || true' EXIT
                curl --fail --retry 20 --retry-delay 1 --retry-connrefused \
                  http://127.0.0.1:18783/ >response.html
                kill "$server_pid"
                wait "$server_pid" || true
                touch "$out"
              '';
        }
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            packages = [
              pkgs.actionlint
              pkgs.go
              pkgs.nixfmt
              pkgs.skopeo
            ];
          };
        }
      );

      formatter = forAllSystems (system: nixpkgs.legacyPackages.${system}.nixfmt);
    };
}
