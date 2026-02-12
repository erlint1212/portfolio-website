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
            echo "Environment: GOTH: Go + Tailwind + Htmx + Templ + K8s"
            echo "------------------------------------------------"
            go version
            templ --version
            alias tgr='templ generate && go run cmd/web/main.go'
            alias k='kubectl'
            alias kgp='kubectl get pods'
            echo "Commands:"
            echo "RabbitMQ: kubectl port-forward service/portfolio-rabbitmq-service 5672:5672"
            echo "Website: kubectl port-forward service/portfolio-service 8000:80"
          '';
        };
      }
    );
}
