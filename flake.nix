{
  description = "mono devshell";
  inputs = {
    #nixpkgs.url = "github:nixos/nixpkgs/nixos-22.11";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }@inputs:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        write-mailmap = pkgs.buildGoModule rec {
          name = "write_mailmap";
          src = pkgs.fetchFromGitHub {
            owner = "CyCoreSystems";
            repo = "write_mailmap";
            rev = "v0.3.0";
            sha256 = "sha256-LzLLEtsWLeIOnlY1pygAOhTsGiWfISnuVF/jeoHHzaw=";
          };

          # There are no upstream packages, so vendor hash is null.
          vendorHash = null;
        };

        gci = pkgs.buildGoModule rec {
          name = "gci";
          src = pkgs.fetchFromGitHub {
            owner = "daixiang0";
            repo = "gci";
            rev = "v0.10.1";
            sha256 = "sha256-/YR61lovuYw+GEeXIgvyPbesz2epmQVmSLWjWwKT4Ag=";
          };

          # Switch to fake vendor sha for upgrades:
          #vendorSha256 = pkgs.lib.fakeSha256;
          vendorSha256 = "sha256-g7htGfU6C2rzfu8hAn6SGr0ZRwB8ZzSf9CgHYmdupE8=";
        };

        cclint = pkgs.writeScriptBin "lint" ''
          cd $(git rev-parse --show-toplevel)
          write_mailmap > CONTRIBUTORS
          gofumpt -w .
          gci write --skip-generated -s standard -s default -s "Prefix(github.com/CyCoreSystems)" .
          golangci-lint run
        '';

        ccmocks = pkgs.writeScriptBin "gen-mocks" ''
           rm -Rf vendor/ client/arimocks
           mockery --name . --outpkg arimocks --output client/arimocks
        '';
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            buf
            cclint
            ccmocks
            gci
            go-tools
            write-mailmap
          ];
        };
      });
}
