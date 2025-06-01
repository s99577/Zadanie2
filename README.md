CI/CD Pipeline zawarty w pliku zadanie2.yml realzuje następujące kroki:

1. Wyzwalacze:
  - workflow_dispatch - umożliwia ręczne uruchomienie workflow przy użyciu GitHub API, GitHub CLI lub GitHub UI.
  - push: branches: -main - automatyczne uruchomienie przy wypchnięciu zmian do gałęzi main
  - push: tags: - 'v*.*.*' - automatyczne uruchomienie przy wypchnięciu tagu w określonym formacie (np. v1.0.0)

2. Zmienne środowiskowye
  - PORT - port aplikacji
  - TZ - strefa czasowa
  - WEATHER_API_KEY - klucz API (pobierany z secrets)
  - TEST_TAG - tag dla obrazu testowego

4. Job ci_step (uruchomiony na ubuntu-24.04):
  - Pobranie kodu źródłowego repozytorium (actions/checkout@v4)
  - Wygenerowanie tagów Dockera na podstawie metadanych Git (docker/metadata-action@v5)
  - Konfiguracja QEMU (docker/setup-qemu-action@v3)
  - Konfiguracja Buildx (docker/setup-buildx-action@v3)
  - Logowanie do DockerHub (danen logowania pobierane z vars i secrets)
  - Logowanie do GitHub (dane logowania pobierane z secrets i github.repository_owner)
  - Budowanie obrazu testowego dla architektury linux/amd64, wykorzystanie danych cache
  - Skanowanie obrazu testowego przy użyciu skanera Trivy (aquasecurity/trivy-action@0.30.0), przy wykryciu zagrożenia sklasyfikowanego jako krytyczne lub wysokie zostanie zwrócony błąd
  - Jeżeli skanowanie nie zwróci błędu nastąpi budowanie obrazu dla architektur linux/amd64 i linux/arm64 oraz przesłany do repozytorium Github

Sposób tagowania obrazów:
  - Dynamiczne:
    1. type=semver,pattern={{version}} - tagowanie semver dla oficjalnych wydań, które umożliwia szybką identyfikację wersji;
    2. type=semver,pattern={{major}} - tagowanie semver dla głównych wersji, który wskazuje najnowszą wersję z danej serii;
    3. type=sha,format=short,prefix=sha_ - tagowanie krótkim hashem commita.
  - Statyczne:
    1. Obraz używany do skanowania CVE jest tymczasowo tagowany statyczną nazwą określoną w env.

Sposób tagowania danych cache:
Dane cache są przechowywane w dedykowanym repozytorium na DockerHub.
Tagowanie jako cache_${{ github.ref_name } pozwala na utrzymanie oddzielengo cache'u dla każdej wersji oraz gałęzi.
Statyczne określenie dodatkowego źródła cache'u jako cache_main pozwala na użycie cache'u gałęzi głównej w przypadku budowania obrazu wersji, która nie ma jeszcze własnego cache'u.
