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
            tailwindcss

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
            alias tgr='templ generate && tailwindcss -i ./internal/assets/css/input.css -o ./assets/css/output.css && go run cmd/web/main.go'
            alias k='kubectl'
            alias tailcomp='tailwindcss -i ./internal/assets/css/input.css -o ./assets/css/output.css'
            alias kgp='kubectl get pods'
            echo "Commands:"
            echo "Tailwind CSS (alias tailcomp): tailwindcss -i ./internal/assets/css/input.css -o ./assets/css/output.css"
            echo "RabbitMQ: kubectl port-forward service/portfolio-rabbitmq-service 5672:5672"
            echo "Website: kubectl port-forward service/portfolio-service 8000:80"
          '';
        };
      }
    );
}
