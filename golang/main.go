package main

import (
	"fmt"
	"sync"
	"math/rand"
	"time"
	"encoding/json"
	"io/ioutil"
    "os"
)


type CoursesJSON struct{
	Courses []CourseJSON `json:"courses"`
}

type CourseJSON struct {
	DateDeDepart 	string 			`json:"DateDeDepart"`
	Ordre 			[]VoitureJSON 	`json:"course"`
}

type VoitureJSON struct {
	Nom 	string `json:"nom"`
	Marque 	string `json:"marque"`
}
var coursesJSON CoursesJSON


type Voiture struct {
	Nom 	string
	Marque 	string
}

type Course struct {
	mu 						sync.Mutex
	Ordre 					[]Voiture
	OrdrePrecedent 			[]Voiture
	dateCoursePrecedente 	string
}
var course Course


func afficherMenuPrincipal()(){
	fmt.Println("1. Ajouter une voiture")
	fmt.Println("2. Afficher la position des coureurs")
	fmt.Println("3. Supprimer une voiture")
	fmt.Println("4. Afficher le nombre d'inscrits à la course")
	fmt.Println("5. Qui est le vainqueur de la dernière course")
	fmt.Println("6. Afficher la dernière course")
	fmt.Println("7. Lancer la course")
	fmt.Println("8. Ouvrir le fichier de sauvegarde JSON")
	fmt.Println("9. Quitter le programme")
}

func ajouterVoiture()(Voiture){
	var voitureAAjouter Voiture
	emptyString := true
	for emptyString {
		emptyString = false
		fmt.Println("\tVeuillez introduire le nom de la voiture")
		fmt.Scanf("%s", &(voitureAAjouter.Nom))
		if voitureAAjouter.Nom == "" || len(voitureAAjouter.Nom) == 0 { 
			fmt.Println("\tErreur, Vous n'avez rien introduit")
			emptyString = true
		}
	}
	emptyString = true
	for emptyString {
		emptyString = false
		fmt.Println("\tVeuillez introduire la marque de la voiture")
		fmt.Scanf("%s", &(voitureAAjouter.Marque))
		if voitureAAjouter.Marque == "" || len(voitureAAjouter.Marque) == 0 { 
			fmt.Println("\tErreur, Vous n'avez rien introduit")
			emptyString = true
		}
	}
	course.Ordre = append(course.Ordre, voitureAAjouter)
	fmt.Println("\tvoiture ", voitureAAjouter.Nom, " ajoutée")
	return voitureAAjouter
}

func checkIfCourseEmpty()(bool){
	return len(course.Ordre) == 0 
}

func afficherCourse()(){
	if checkIfCourseEmpty() {
		fmt.Println("\tErreur, aucun coureur inscrit dans la course")
	} else {
		for i, val := range course.Ordre {
			fmt.Printf("\t%d: %s\n", i, val)
		}
	}
}

func supprimerVoiture()(){
	if checkIfCourseEmpty() {
		fmt.Println("\tErreur, aucun coureur inscrit dans la course")
	} else {
		var choixInt int
		doItAgain := true
		for doItAgain {
			fmt.Println("\tQuel courreur voulez vous enlever ? (donnez sa position)")
			for i, val := range course.Ordre {
				fmt.Printf("\t%d: %s\n", i, val)
			}
			
			_, err := fmt.Scanf("%d", &choixInt)
			if err != nil {
				fmt.Println("\tErreur, vous n'avez pas introduit un chiffre")
				doItAgain = true
				continue
			}
			if choixInt < 0 || choixInt > len(course.Ordre)-1 {
				fmt.Println("\tErreur, le chiffre introduit ne correspond pas à une position valide")
				doItAgain = true
				continue
			}
			doItAgain = false
		}
		copy(course.Ordre[choixInt:], course.Ordre[choixInt+1:]) // Shift a[i+1:] left one index.
		course.Ordre[len(course.Ordre)-1] = (Voiture{})     				 // Erase last element (write zero value).
		course.Ordre = course.Ordre[:len(course.Ordre)-1]     	 // Truncate slice.
		fmt.Println("\tSupression terminée, coureurs restants :")
		afficherCourse()
	}
}

func afficherNbCoureursInscrits()(){
	fmt.Println("\tIl y a actuellement ", len(course.Ordre), " inscrits")
}

func checkIfCoursePrecedenteEmpty()(bool){
	return len(course.OrdrePrecedent) == 0 
}

func vainqueurDerniereCourse()(){
	if checkIfCoursePrecedenteEmpty() {
		fmt.Println("\tErreur, aucun coureur inscrit dans la course précédente")
	} else {
		fmt.Println("\tLe vainqueur était ", course.OrdrePrecedent[0])
	}
}

func afficherDerniereCourse()(){
	if checkIfCoursePrecedenteEmpty() {
		fmt.Println("\tErreur, aucun coureur inscrit dans la course précédente")
	} else {
		for i, val := range course.OrdrePrecedent {
			fmt.Printf("\t%d: %s\n", i, val)
		}
	}
}

func lancerCourse()(){
	if len(course.Ordre) < 2 {
		fmt.Println("\tErreur, il n'y a pas assez de coureurs pour lancer la course")
	} else {
		nbrTotalCoureurs := len(course.Ordre)
		fmt.Println("Il y a ", nbrTotalCoureurs, " coureurs")
		course.OrdrePrecedent = make([]Voiture, nbrTotalCoureurs)
		copy(course.OrdrePrecedent, course.Ordre)

		var wg sync.WaitGroup
		course.Ordre = make([]Voiture, 0)
		goDeplacement := func(coureur int) {
			min := 1
			max := nbrTotalCoureurs
			tempsDeParcours := rand.Intn(max - min) + min
			time.Sleep(time.Duration(tempsDeParcours) * time.Second)
			course.mu.Lock()
			defer course.mu.Unlock()
			// zone critique
			course.Ordre = append(course.Ordre, course.OrdrePrecedent[coureur])
			fmt.Println("\tcoureur ", course.OrdrePrecedent[coureur].Nom, " arrivé en ", len(course.Ordre), " position")
			wg.Done()
		}
		// 1 goroutine pour chaque coureur
		wg.Add(nbrTotalCoureurs)
		course.dateCoursePrecedente = time.Now().Format("02-01-2006_15:04:05")
		for i, _ := range course.OrdrePrecedent {
			go goDeplacement(i)
		}
		wg.Wait()
		fmt.Println("Course terminée")
		var courseJson CourseJSON
		for i, val := range course.Ordre {
			fmt.Printf("\t%d: %s\n", i, val)
			voiture := VoitureJSON{course.Ordre[i].Nom, course.Ordre[i].Marque}
			courseJson.Ordre = append(courseJson.Ordre, voiture)
		}
		// sauvegarder la course
		copy(course.OrdrePrecedent, course.Ordre)
		// reset à 0
		course.Ordre = make([]Voiture, 0)
	}
}

func afficherMenuJSON()(){
	fmt.Println("1. Sauvegarder la course précédente dans le fichier JSON")
	fmt.Println("2. Afficher les courses précédentes")
	fmt.Println("3. Supprimer une course précédente")
	fmt.Println("4. Supprimer toutes les courses précédentes")
	fmt.Println("5. Revenir au menu principal")
}

func lectureFichier()(){
	// récupérer les courses déjà inscrits dans le fichier
	jsonFile, err := os.Open("courses.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &coursesJSON)
}

func checkFileEmpty()(bool){
	lectureFichier()
	return len(coursesJSON.Courses) < 1
}

func sauvegarderContenuJSON()(){
	if checkIfCoursePrecedenteEmpty() {
		fmt.Println("\tErreur, aucun coureur inscrit dans la course précédente")
	} else {
		lectureFichier()
		// ajouter en fin, la nouvelle course à sauvegarder
		var courseJson CourseJSON
		for i, _ := range course.OrdrePrecedent {
			voiture := VoitureJSON{course.OrdrePrecedent[i].Nom, course.OrdrePrecedent[i].Marque}
			courseJson.Ordre = append(courseJson.Ordre, voiture)
		} 
		courseJson.DateDeDepart = course.dateCoursePrecedente
		coursesJSON.Courses = append(coursesJSON.Courses, courseJson)
		file, _ := json.MarshalIndent(coursesJSON, "", "   ")
		_ = ioutil.WriteFile("courses.json", file, 0644)
		fmt.Println("\tCourse sauvegardée")
	}
}

func afficherContenuJSON()(){
	if checkFileEmpty() {
		fmt.Println("\tErreur, il n'y a aucune course inscrite dans le fichier JSON")
	} else {
		for i := 0; i < len(coursesJSON.Courses); i++ {
			fmt.Println("date de la course : " + coursesJSON.Courses[i].DateDeDepart)
			for j, _ := range coursesJSON.Courses[i].Ordre {
				fmt.Println("\t{ Nom : ", coursesJSON.Courses[i].Ordre[j].Nom, " -- Marque : ", coursesJSON.Courses[i].Ordre[j].Marque, "}")
			}
		}
	}
}

func supprimerCourse()(){
	if checkFileEmpty() {
		fmt.Println("\tErreur, il n'y a aucune course inscrite dans le fichier JSON")
	} else {
		for i := 0; i < len(coursesJSON.Courses); i++ {
			fmt.Println("Indice (", (i+1), ") - date de la course : " + coursesJSON.Courses[i].DateDeDepart)
			for j, _ := range coursesJSON.Courses[i].Ordre {
				fmt.Println("\t{ Nom : ", coursesJSON.Courses[i].Ordre[j].Nom, " -- Marque : ", coursesJSON.Courses[i].Ordre[j].Marque, "}")
			}
		}
		fmt.Println("Veuillez choisir l'indice de la course à supprimer")
		var choixInt int
		for choixInt < 1 || choixInt > len(coursesJSON.Courses) {
			fmt.Println("veuillez introduire votre choix :")
			_, err := fmt.Scanf("%d", &choixInt)
			if err != nil {
				fmt.Println("\tErreur, vous n'avez pas introduit un chiffre")
			}
			if choixInt < 1 || choixInt > len(coursesJSON.Courses) {
				fmt.Println("\tErreur, le chiffre introduit ne correspond pas à un choix valide")
			} else {
				break
			}
			choixInt = 0
			fmt.Println()
		}
		choixInt -= 1
		copy(coursesJSON.Courses[choixInt:], coursesJSON.Courses[choixInt+1:])
		coursesJSON.Courses[len(coursesJSON.Courses)-1] = CourseJSON{}
		coursesJSON.Courses = coursesJSON.Courses[:len(coursesJSON.Courses)-1] 
		file, _ := json.MarshalIndent(coursesJSON, "", "   ")
		_ = ioutil.WriteFile("courses.json", file, 0644)
		fmt.Println("Suppression effectuée")
	}
}

func supprimerCourses()(){
	if checkFileEmpty() {
		fmt.Println("\tErreur, il n'y a aucune course inscrite dans le fichier JSON")
	} else {
		coursesJSON.Courses = make([]CourseJSON, 0)
		file, _ := json.MarshalIndent(coursesJSON, "", "   ")
		_ = ioutil.WriteFile("courses.json", file, 0644)
		fmt.Println("Suppression effectuée")
	}
}

func menuFichierJSON()(){
	var choixInt int
	for choixInt < 1 || choixInt > 5 {
		afficherMenuJSON()
		fmt.Println("veuillez introduire votre choix :")
		_, err := fmt.Scanf("%d", &choixInt)
		if err != nil {
			fmt.Println("Erreur, vous n'avez pas introduit un chiffre")
		}
		if choixInt < 1 || choixInt > 9 {
			fmt.Println("Erreur, le chiffre introduit ne correspond pas à un choix valide")
		}
		switch choixInt {
			case 1:
				sauvegarderContenuJSON()
			case 2:
				afficherContenuJSON()
			case 3:
				supprimerCourse()
			case 4:
				supprimerCourses()
			case 5:
				return
		}
		choixInt = 0
		fmt.Println()
	}
}

func main(){
	fmt.Println("Bienvenue dans le programme de course")
	var choixInt int
	for choixInt < 1 || choixInt > 9 {
		afficherMenuPrincipal()
		fmt.Println("veuillez introduire votre choix :")
		_, err := fmt.Scanf("%d", &choixInt)
		if err != nil {
			fmt.Println("Erreur, vous n'avez pas introduit un chiffre")
		}
		if choixInt < 1 || choixInt > 9 {
			fmt.Println("Erreur, le chiffre introduit ne correspond pas à un choix valide")
		}
		switch choixInt {
			case 1:
				voitureAAjouter := ajouterVoiture()
				if voitureAAjouter == (Voiture{}) {
					fmt.Println("\tErreur d'ajout")
				} else {
					fmt.Println("\tLe nom -> ", voitureAAjouter.Nom,"\n\tLa marque -> ", voitureAAjouter.Marque)
				}
			case 2:
				afficherCourse()
			case 3:
				supprimerVoiture()
			case 4:
				afficherNbCoureursInscrits()
			case 5:
				vainqueurDerniereCourse()
			case 6:
				afficherDerniereCourse()
			case 7:
				lancerCourse()
			case 8:
				menuFichierJSON()
			case 9:
				fmt.Println("\tMerci, aurevoir")
				return
		}
		choixInt = 0
		fmt.Println()
	}
}
