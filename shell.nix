let
  pkgs = import <nixpkgs> { config.allowUnfree = true; };
in pkgs.mkShell {
  name = "go-htmx-k8s-dev-shell";

  packages = with pkgs; [
    # Core Languages & Frameworks
    go              # The Go programming language
    templ           # The templating engine for Go

    # Infrastructure & Containerization
    docker          # Docker client
    kubectl         # Kubernetes command-line tool
    
    # Messaging
    rabbitmq-server # RabbitMQ server for local development
  ];

  shellHook = ''
    echo "Welcome to your Go + Htmx + K8s dev environment!"
    
    # Set GOPATH to a local directory if you prefer not to use the global one
    # export GOPATH=$PWD/.gopath
    # export PATH=$GOPATH/bin:$PATH

    # Display versions for verification
    echo "Tools loaded:"
    go version
    templ --version
    kubectl version --client
    echo "------------------------------------------------"
  '';
}
