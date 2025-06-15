# Symulacja Lisów i Królików

## Opis projektu
Projekt przedstawia symulację ekosystemu, w którym lisy oraz króliki współistnieją w środowisku z odrastającą trawą.

## Funkcjonalności
- Interaktywna symulacja z wizualizacją graficzną
- Lisy polują na króliki dla energii
- Króliki jedzą trawę i uciekają przed lisami. Mogą "potknąć" się uciekając
- Zwierzęta mogą się rozmnażać gdy spełnione są odpowiednie warunki
  - Warunki rozmnażania lisów:
    - Energia lisa musi być wyższa niż koszt reprodukcji
    - W pobliżu musi znajdować się inny lis
    - Od poprzedniej reprodukcji musi upłynąć określona ilość tur
  - Warunki rozmnażania królików:
    - Energia królika musi być wyższa niż koszt reprodukcji
    - W pobliżu musi znajdować się inny królik
    - Od poprzedniej reprodukcji musi upłynąć określona liczba tur
- Trawa odrasta z czasem
  - Gdy trawa zostanie zjedzona, rozpoczyna się odliczanie do jej odrośnięcia
  - Trawa odrasta do określonej ilości w jednym polu

## Instalacja
```bash
go mod init
```

## Uruchomienie
```bash
go run main.go
```

## Sterowanie
- Lewy przycisk myszy: Dodaj królika
- Prawy przycisk myszy: Dodaj lisa
- Można "rysować"

## Konfiguracja
Parametry symulacji można zmieniać w pliku config.go.

## Ciekawe spostrzeżenia

- Intrygująca zależność: króliki nie mogą żyć z lisami, ale lisy nie mogą żyć bez królików.

- Presja drapieżników powoduje migrację królików. Mimo pewnych ofiar, śledzenie ich przez lisy nieraz potrafi posłużyć populacji królików. Powoduje to uciekanie królików w dalsze miejsca, co może pozwolić na uniknięcie krwawej rzezi w poprzednim miejscu bytowania. 

- Bardzo duż wpływ na stabilność populacji obydwu gatunków było wprowadzenie "potknięć" podczas ucieczki. Królik uciekając przed lisem ma 20% szans na zostanie w miejscu.

- Problematyczne mnożenie się królików "jak króliki" zostało unormowane przez ograniczenie dostępności trawy: maksymalnie 2 w jednym miejscu z długim czasem odrastania.

- Przy aktualnych ustawieniach, udało mi się doprowadzić do dobrej stabilizacji symulacji. W większości wypadków symulacja trwa bez końca z obydwoma gatunkami.