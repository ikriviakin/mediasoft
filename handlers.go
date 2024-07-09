package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
)

// Функция для чтения списка автомобилей из файла
func ReadCarsFromFile() ([]Car, error) {
	file, err := os.Open("cars.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cars []Car
	err = json.NewDecoder(file).Decode(&cars)
	if err != nil {
		return nil, err
	}
	return cars, nil
}

// Функция для записи списка автомобилей в файл
func WriteCarsToFile(cars []Car) error {
	file, err := os.Create("cars.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(cars)
}

// Функция для генерации нового уникального ID
func GenerateNewID(cars []Car) string {
	maxID := 0
	for _, car := range cars {
		id, err := strconv.Atoi(car.ID)
		if err == nil && id > maxID {
			maxID = id
		}
	}
	return strconv.Itoa(maxID + 1)
}

// Обработчик для создания нового автомобиля
func CreateEntity(w http.ResponseWriter, r *http.Request) {
	var newCar Car
	err := json.NewDecoder(r.Body).Decode(&newCar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cars, err := ReadCarsFromFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверка, предоставлен ли ID
	if newCar.ID == "" {
		newCar.ID = GenerateNewID(cars)
	}

	cars = append(cars, newCar)

	err = WriteCarsToFile(cars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCar)
}

// Обработчик для обновления всех автомобилей
func UpdateEntities(w http.ResponseWriter, r *http.Request) {
	var updatedCars []Car
	err := json.NewDecoder(r.Body).Decode(&updatedCars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = WriteCarsToFile(updatedCars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCars)
}

// Обработчик для получения списка всех автомобилей
func GetCars(w http.ResponseWriter, r *http.Request) {
	cars, err := ReadCarsFromFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

// Обработчик для получения одного автомобиля по ID
func GetCar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	cars, err := ReadCarsFromFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, car := range cars {
		if car.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(car)
			return
		}
	}

	http.Error(w, "Car not found", http.StatusNotFound)
}

// Обработчик для обновления информации об автомобиле
func UpdateCar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var updatedCar Car
	err := json.NewDecoder(r.Body).Decode(&updatedCar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cars, err := ReadCarsFromFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, car := range cars {
		if car.ID == id {
			if updatedCar.ID == "" {
				updatedCar.ID = car.ID // Сохранение старого ID, если новый не указан
			}
			cars[i] = updatedCar
			err = WriteCarsToFile(cars)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedCar)
			return
		}
	}

	http.Error(w, "Car not found", http.StatusNotFound)
}

// Обработчик для частичного обновления информации об автомобиле
func PatchCar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var patchData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&patchData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cars, err := ReadCarsFromFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, car := range cars {
		if car.ID == id {
			if brand, ok := patchData["brand"].(string); ok {
				car.Brand = brand
			}
			if model, ok := patchData["model"].(string); ok {
				car.Model = model
			}
			if mileage, ok := patchData["mileage"].(float64); ok {
				car.Mileage = int(mileage)
			}
			if ownersCount, ok := patchData["owners_count"].(float64); ok {
				car.OwnersCount = int(ownersCount)
			}

			cars[i] = car

			err = WriteCarsToFile(cars)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent) // Устанавливаем статус 204 No Content
			return
		}
	}

	http.Error(w, "Car not found", http.StatusNotFound)
}

// Обработчик для удаления автомобиля по ID
func DeleteCar(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	cars, err := ReadCarsFromFile()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for i, car := range cars {
		if car.ID == id {
			cars = append(cars[:i], cars[i+1:]...)

			err = WriteCarsToFile(cars)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent) // Устанавливаем статус 204 No Content
			return
		}
	}

	http.Error(w, "Car not found", http.StatusNotFound)
}
