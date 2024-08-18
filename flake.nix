{
  description = "Nix flake for Neon";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      gomod2nix,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ gomod2nix.overlays.default ];
        };
        readFileOr = (path: default: with builtins; if pathExists path then (readFile path) else default);
        repoName = "github.com/bow/neon";
        tagFile = "${self}/.tag";
        revFile = "${self}/.rev";
        version = readFileOr tagFile "0.0.0";
        commit = readFileOr revFile "";
        app = gomod2nix.legacyPackages.${system}.buildGoApplication {
          src = ./.;
          pwd = ./.;
          name = "neon";
          doCheck = false;
          CGO_ENABLED = 0;
          ldflags = [
            "-w" # do not generate debug output
            "-s" # strip symbols table
            "-X ${repoName}/internal.version=${version}"
            "-X ${repoName}/internal.gitCommit=${commit}"
          ];
        };
      in
      {
        devShells = {
          default = pkgs.mkShell {
            packages = with pkgs; [
              # go-only
              go
              gocover-cobertura
              golangci-lint
              gomod2nix.packages.${system}.default
              gopls
              gosec
              gotestsum
              gotools
              (go-migrate.overrideAttrs (_final: _prev: { tags = [ "sqlite" ]; }))
              mockgen
              protobuf
              protoc-gen-go
              protoc-gen-go-grpc
              # nix-only
              deadnix
              nixfmt-rfc-style
              statix
              # others
              gnugrep
              sqlite
            ];
          };
        };
        packages =
          let
            imgTag = readFileOr tagFile "latest";
            imgAttrs = rec {
              name = "ghcr.io/bow/${app.name}";
              tag = imgTag;
              contents = [ app ];
              config = {
                Entrypoint = [ "/bin/${app.name}" ];
                Env = [
                  "NEON_SERVER_ADDR=tcp://0.0.0.0:5151"
                  "NEON_SERVER_DB_PATH=/var/data/neon.db"
                ];
                Labels = {
                  "org.opencontainers.image.revision" = readFileOr revFile "";
                  "org.opencontainers.image.source" = "https://github.com/bow/${app.name}";
                  "org.opencontainers.image.title" = "${app.name}";
                  "org.opencontainers.image.url" = "https://${name}";
                };
              };
              extraCommands = ''
                mkdir -p var/data
              '';
            };
          in
          {
            dockerArchive = pkgs.dockerTools.buildLayeredImage imgAttrs;
            dockerArchiveStreamer = pkgs.dockerTools.streamLayeredImage imgAttrs;
            local = app;
          };
      }
    );
}
