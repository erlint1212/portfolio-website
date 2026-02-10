{
  description = "Backend Archmage Portfolio Environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-25.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            templ
            docker
            kubectl
            minikube
            rabbitmq-server

            python313Packages.pathspec
            
            # Useful tools for debugging
            curl 
            jq
          ];

          shellHook = ''
            echo "Enter the Minikube Docker Environment"
            eval $(minikube docker-env)
            echo "Environment: Go + Htmx + Templ + K8s"
            echo "------------------------------------------------"
            go version
            templ --version
            alias tgr='templ generate && go run cmd/web/main.go'
          '';
        };
      }
    );
}
