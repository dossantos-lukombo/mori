
graph TD
  subgraph 🛠️ Développement
    A[Commit sur branche feature/*] -->|Pull Request| B[PR vers develop]
    B --> C{CI : lint + tests<br/>build docker}
    C -->|Succès| D[Merge sur develop]
  end

  D --> E[semantic-release<br/>🔀 bump version<br/>🏷️ tag vX.Y.Z]
  E --> F[PR « develop → production »]

  subgraph 🚀 Production
    F --> G{CI/CD}
    G --> H[Build images Docker]
    H --> I[Push vers Docker Hub]
    I --> J[ssh + docker-compose up]
    J --> K[💚 Site en ligne, café chaud]
  end
