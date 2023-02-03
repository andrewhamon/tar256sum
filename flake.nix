{
  description = "Packages and modules for tar256sum";

  inputs.nixpkgs.url = "nixpkgs/nixos-22.11";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  inputs.nixpkgs-old-git.url = "nixpkgs/nixos-22.05";

  outputs = { self, nixpkgs, flake-utils, nixpkgs-old-git }: flake-utils.lib.eachDefaultSystem
    (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        oldgit = nixpkgs-old-git.legacyPackages.${system}.git;
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "tar256sum";
          version = "0.0.1";
          vendorSha256 = "sha256-pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";
          src = ./.;
        };
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = [
            pkgs.go
            pkgs.gopls
            pkgs.nixpkgs-fmt
            pkgs.delve
            pkgs.go-tools
            pkgs.git
            (pkgs.writeScriptBin "oldgit" ''
              ${oldgit}/bin/git $@
            '')
          ];
        };
      }
    );
}
