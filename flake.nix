{
  description = "Backend Portfolio Website Environment";

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
            #  Core Stack 
            go
            templ
            tailwindcss_4

            #  Kubernetes 
            docker
            kubectl
            k9s
            minikube
            kubernetes-helm     # Deploying ArgoCD + Vaultwarden charts

            #  CI/CD 
            argocd              # ArgoCD CLI (sync, app management)

            #  Secrets & TLS 
            openssl             # Generate self-signed certs
            mkcert              # Local-trusted dev certs (no browser warnings)
            nss.tools           # certutil — needed by mkcert for browser CA install

            #  Messaging 
            rabbitmq-server

            #  Remote / Mount 
            sshfs

            #  Build Context 
            python313Packages.pathspec
            
            #  Debugging 
            curl 
            jq
          ];

          shellHook = ''
            echo "Enter the Minikube Docker Environment"
            eval $(minikube docker-env)
            echo "  GOTH Stack + server: Go + Tailwind + Htmx + Templ"
            echo "  + K8s + ArgoCD + Vaultwarden"
            go version
            templ --version
            echo ""

            #  Aliases 
            alias k='kubectl'
            alias kgp='kubectl get pods'
            alias kga='kubectl get all'
            alias tgr='templ generate && tailwindcss -i ./internal/assets/css/input.css -o ./internal/views/css/output.css && go run cmd/web/main.go'
            alias tailcomp='tailwindcss -i ./internal/assets/css/input.css -o ./internal/views/css/output.css'

            echo "Commands:"
            echo "  tgr        — templ generate + tailwind + go run"
            echo "  tailcomp   — recompile Tailwind CSS"
            echo "  kgp        — kubectl get pods"
            echo ""
            echo "Port forwards:"
            echo "  RabbitMQ:  kubectl port-forward svc/portfolio-rabbitmq-service 5672:5672"
            echo "  Website:   kubectl port-forward svc/portfolio-service 8000:80"
            echo "  ArgoCD:    kubectl port-forward svc/argocd-server -n argocd 8083:443"
            echo "  Vaultwarden: kubectl port-forward svc/vaultwarden -n vaultwarden 8443:80"
          '';
        };
      }
    );
}

