{
  description = "A basic gomod2nix flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/9355fa86e6f27422963132c2c9aeedb0fb963d93";
    flake-utils.url = "github:numtide/flake-utils/b1d9ab70662946ef0850d488da1c9019f3a9752a";
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
        app = gomod2nix.legacyPackages.${system}.buildGoApplication {
          src = ./.;
          pwd = ./.;
          name = "neon";
          CGO_ENABLED = 0;
          ldflags = [
            "-w" # do not generate debug output
            "-s" # strip symbols table
          ];
        };
      in
      {
        packages = {
          default = app;
        };
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
              pre-commit
              sqlite
            ];
          };
        };
      }
    );
}
