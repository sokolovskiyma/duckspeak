package main

// iconv -f cp1251 -t utf8 Grin_A._Alyie_Parusa.txt -o 27764733.txt

import (
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

// Declares Dictionary // Обявляет Словарь
var Dictionary = map[string]map[string]int{}

// File name // Имя файла
var DicName = "dictionary.gob"

// Loads the dictionary from disk, // Загружает Словарь с диска,
// or if it is empty, sets the initial values // или если он пустой проставляет начальные значения
func init() {
	// If the file is UNAVALIBLE, initializes the initial value of the Dictionary // Если файл НЕДОСТУПЕН, инициализирует начальное значение Словаря
	if _, err := os.Stat(DicName); os.IsNotExist(err) {
		Dictionary["*START*"] = map[string]int{}
		// If the file is available, loads data from it into the Dictionary // Если файл доступен, загружает данные из него в Словарь
	} else {
		// Read the file // Читает файл
		dataFile, err := os.Open(DicName)
		if err != nil {
			panic(err)
		}
		defer dataFile.Close()

		// Decodes binary data // Декодирует бинарные данные
		dataDecoder := gob.NewDecoder(dataFile)
		err = dataDecoder.Decode(&Dictionary)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	// Rand seed
	rand.Seed(time.Now().Unix())
	// Declares and parses flags // Объявляет и парсит флаги
	learn := flag.String("l", "", "single file or folder to add to the dictionary")
	generate := flag.Int("g", 0, "generate n random sentences")
	clear := flag.Bool("c", false, "Clear the dictionary data")
	flag.Parse()
	// If there is a flag "learn" - adds values from the text to the Dictionary // Если есть флаг "learn" - добавляет значения из текста в Словарь
	if *learn != "" {
		learnFiles(*learn)
	}
	// If the argument flag "generate" > 0 // Если аргумент флага "generate" > 0
	// generates" generate " random sentences // генерирует "generate" случайных предложений
	if *generate > 0 {
		// If the dictionary has not been trained, notify the // Если словарь не обучен, сообщает об этом
		if len(Dictionary) < 2 {
			fmt.Println("Teach with -l flag first")
			// Otherwise, generates the "generate" random sentences // Иначе генерирует "generate" случайных предложений
		} else {
			for i := 0; i < *generate; i++ {
				fmt.Print(generateSentence() + " ")
			}
		}
		fmt.Print("\n")
	}

	// If the "clear" flag is set - clears the Dictionary  // Если установлен флаг "clear" - очищает Словарь
	if *clear {
		clearDictionary()
	}
}

// Reads files / folders // Читает файлы / папки
// Sends files one at a time to the add to Dictionary function // Отправляет файлы по одному в функцию addToDictionary
func learnFiles(path string) {
	// Gets information on files // Получает информацию по файлам
	fileInfo, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	// file mode bits
	fileMode := fileInfo.Mode()
	// If it's a folder // Если это папка
	if fileMode.IsDir() {
		// Reads folder // Читает папку
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		// For each item, check that it is not a folder // Для каждого элемента, проверяет что это не папка
		for _, file := range files {
			if !file.IsDir() {
				// Reads each file // Читает каждый файл
				data, err := ioutil.ReadFile(path + file.Name())
				if err != nil {
					panic(err)
				}
				// Sends a file to the Dictionary // Отправляет файл в Словарь
				fmt.Println(file.Name())
				addToDictionary(string(data))
			}

		}
		// If this is a regular file // Если это регулярный файл
	} else if fileMode.IsRegular() {
		// Read it // Читает его
		data, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		// Sends a file to the Dictionary // Отправляет файл в Словарь
		fmt.Println(fileInfo.Name())
		addToDictionary(string(data))
	}
}

// Parses the text with regexps adds a pair key-value in the Dictionary // Парсит текст регулярками добавляет пары ключ-значение в Словарь
func addToDictionary(text string) {
	// At the end saves the Dictionary // В конце сохраняет Словарь
	defer saveDictionary()
	// Word count in text // Счетчик слов в тексте
	var wordCount int
	// Regex // Регулярки
	onlyChar := regexp.MustCompile(`[^A-zА-я ]`)
	doubleSpace := regexp.MustCompile(` +`)
	// Splits text into sentences // Разбивает текст на предложения
	sentences := regexp.MustCompile(`[\.\?\!\:\;…]{1,3}\s`).Split(text, -1)
	for _, sentence := range sentences {
		// Clears offers from !letters // Очищает предложения от !букв
		sentence = onlyChar.ReplaceAllString(sentence, " ")
		sentence = doubleSpace.ReplaceAllString(sentence, " ")
		sentence = strings.TrimSpace(sentence)
		// Ignores empty strings // Игнорирует пустые строки
		if sentence != "" {
			// Splits the string into separate words according to the delimiter " " // Разбивает строки на отдельные слова по разделителю " "
			words := strings.Split(sentence, " ")
			// Lowercase all words // Приводит все слова к нижнему регистру
			for index, word := range words {
				words[index] = strings.ToLower(word)
			}
			// For each word checks if it is in the Dictionary // Для каждого слова проверяет есть оно в Словаре
			for index, word := range words {
				// If the word can be the beginning of a sentence then the key to it will be *START* // Если слово может быть в начале предложения то ключом к нему будет *START*
				if index == 0 {
					if _, avalable := Dictionary["*START*"][word]; avalable {
						Dictionary["*START*"][word]++
					} else {
						Dictionary["*START*"][word] = 1
					}
				} else {
					// If the previous word is not in the dictionary, // Если предыдущее слово не в словаре,
					// adds it and initializes it with an empty map[string]int value{} // добавляет его и инициализирует пустым значением map[string]int{}
					if _, avalable := Dictionary[words[index-1]]; !avalable {
						Dictionary[words[index-1]] = map[string]int{}
					}
					// If the word is already in the values of the previous, // Если слово уже есть в значениях предыдущего,
					// increases the number of occurrences by one // увеличивает количество вхождений на единицу
					if _, avalable := Dictionary[words[index-1]][word]; avalable {
						Dictionary[words[index-1]][word]++
						// If the word is not in the values of the previous, // Если слова нет в значениях предыдущего,
						// adds it with the number of occurrences = 1 // добавляет его с количеством вхождений = 1
					} else {
						Dictionary[words[index-1]][word] = 1
					}
				}

				// If this is the last word adds to it the key *END* // Если это последнее слово добавляет к немк ключ *END*
				if index == len(words)-1 {
					if _, avalable := Dictionary[words[index]]; !avalable {
						Dictionary[words[index]] = map[string]int{}
					}
					if _, avalable := Dictionary[words[index]]["*END*"]; avalable {
						Dictionary[words[index]]["*END*"]++
					} else {
						Dictionary[words[index]]["*END*"] = 1
					}
				}
				// Increases the count of the total number of words // Увеличивает счетчик общего количества слов
				wordCount++
				fmt.Printf("\r%d words added", wordCount)
			}
		}
	}
	fmt.Print("\n")
}

// Generates a random sentence // Генерирует случайное предложение
func generateSentence() string {
	var resultArr []string
	// Starts with the word " * START*", which is the beginning of the sentence // Начинает со слова "*START*", что является началом предложения
	// at each iteration, assigns the currentWord added to the word clause // на каждой итерации присваивает currentWord добавленное в предложение слово
	for currentWord := "*START*"; currentWord != "*END*"; {
		var temp []string
		// For all words that are values of the currentWord key // Для всех слов, которые являются значениями ключа currentWord
		for key, value := range Dictionary[currentWord] {
			// Adds a word-value to the temporary array as many times // Добавляет во временный массив слово-значение столько раз,
			// as it is a currentWord value // сколько раз оно является значением currentWord
			for i := 0; i < value; i++ {
				temp = append(temp, key)
			}
		}
		// Selects a random word from a temporary array, // Выбирает случайное слово из временного массива,
		// adds it to the resulting array // добавляет его в итоговый массив
		currentWord = temp[rand.Intn(len(temp))]
		resultArr = append(resultArr, currentWord)
	}
	// Makes the first letter uppercase, cuts off " *END*", adds a dot to the end // Делает первую букву заглавной, отсекает "*END*", добавляет точку в конец
	resultArr[0] = strings.Title(resultArr[0])
	resultArr = resultArr[:len(resultArr)-1]
	resultArr[len(resultArr)-1] += "."
	// Returns all words in the array joined to a string with a separator " " // Возвращяет все слова из массива, соединенные в строку, с разделителем " "
	return strings.Join(resultArr, " ")
}

// Saves Dictionary To file // Сохраняет Словарь в файл
func saveDictionary() {
	// Creates a file, if it already exists - opens it // Создает файл, если он уже есть - открывает его
	dataFile, err := os.Create(DicName)
	defer dataFile.Close()
	if err != nil {
		panic(err)
	}
	// Encodes a Dictionary into a binary data, writes to file // Кодирует Словарь в бинарные данные, записывает в файл
	dataEncoder := gob.NewEncoder(dataFile)
	err = dataEncoder.Encode(Dictionary)
	if err != nil {
		panic(err)
	}
}

// Deletes the file with the Dictionary // Удаляет файл со Словарем
func clearDictionary() {
	// If the file is UNAVAILABLE - does nothing // Если файл НЕДОСТУПЕН - не делает ничего
	if _, err := os.Stat(DicName); os.IsNotExist(err) {
		fmt.Println("Nothing to do")
		// If the file is available, deletes it // Если файл доступен удаляет его
	} else {
		err := os.Remove(DicName)
		if err != nil {
			panic(err)
		}
		fmt.Println("Dictionary cleaned")
	}
}
