package main

import triple_s "github.com/ab-dauletkhan/triple-s/cmd/triple-s"

func main() {
	triple_s.Run()
	// fileName := "./data/buckets.xml"

	// // Step 1: Try to read the XML file
	// xmlData, err := os.ReadFile(fileName)
	// if err != nil {
	// 	// Step 2: Check if the file does not exist
	// 	if errors.Is(err, os.ErrNotExist) {
	// 		fmt.Println("File does not exist, creating a new one...")

	// 		// Create a new bucket (or multiple buckets if you want)
	// 		newBucket := Bucket{
	// 			Name:       "NewBucket",
	// 			CreatedAt:  time.Now(),
	// 			ModifiedAt: time.Now(),
	// 			Status:     "New",
	// 		}

	// 		newBucket2 := Bucket{
	// 			Name:       "NewBucket2",
	// 			CreatedAt:  time.Now(),
	// 			ModifiedAt: time.Now(),
	// 			Status:     "New",
	// 		}

	// 		// Prepare the initial data
	// 		buckets := []Bucket{newBucket, newBucket2}

	// 		// Marshal the new bucket data to XML
	// 		newXML, err := xml.MarshalIndent(buckets, "", "  ")
	// 		if err != nil {
	// 			fmt.Println("Error marshaling new XML:", err)
	// 			return
	// 		}

	// 		// Ensure the directory exists before writing the file
	// 		err = os.MkdirAll("./data", 0o755)
	// 		if err != nil {
	// 			fmt.Println("Error creating directory:", err)
	// 			return
	// 		}

	// 		// Write the new XML data to the file
	// 		err = os.WriteFile(fileName, newXML, 0o644)
	// 		if err != nil {
	// 			fmt.Println("Error creating new file:", err)
	// 			return
	// 		}

	// 		fmt.Println("New file created with default data!")
	// 		return
	// 	} else {
	// 		// Other errors like permission issues
	// 		fmt.Println("Error reading file:", err)
	// 		return
	// 	}
	// }

	// // Step 3: File exists, so unmarshal the XML data
	// var buckets []Bucket
	// err = xml.Unmarshal(xmlData, &buckets)
	// if err != nil {
	// 	fmt.Println("Error unmarshaling XML:", err)
	// 	return
	// }

	// // Step 4: Modify the data
	// for i := range buckets {
	// 	// Example: Update the Status and ModifiedAt fields
	// 	buckets[i].Status = "Updated"
	// 	buckets[i].ModifiedAt = time.Now()
	// }

	// // Step 5: Marshal the updated data back to XML
	// updatedXML, err := xml.MarshalIndent(buckets, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error marshaling updated XML:", err)
	// 	return
	// }

	// // Step 6: Write the updated XML back to the file
	// err = os.WriteFile(fileName, updatedXML, 0o644)
	// if err != nil {
	// 	fmt.Println("Error writing updated file:", err)
	// 	return
	// }

	// fmt.Println("File updated successfully!")
}
