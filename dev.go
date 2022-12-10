// func getDebet() {
// 	var dbName string = "expenses.db"
//      //  Initialise DB
// 	db, err := sql.Open("sqlite3", dbName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()
// 	err = db.Ping()
//
// 	if err != nil {
// 		log.Fatalf("ping failed: %s", err)
//     }
//
// //     stmt, err = db.Prepare("SELECT debetCredit from expenses WHERE ")
// //     if err != nil {
// //         log.Fatalf("insert prepare failed: %s", err)
// //     }
//
// //     _, err = stmt.Exec("SELECT * from expenses WHERE debetCredit = 'Debet'")
// //     if err != nil {
// //         log.Fatalf("insert failed(%s, %s, %s): %s", boekdatum, rekeningnummer, bedrag, debetCredit, err)
// //     }
//
// 	stmt, err := db.Prepare("SELECT * from expenses WHERE debetCredit = 'Debet'")
// 	if err != nil {
// 		log.Fatalf("prepare failed: %s", err)
// 	}
//
// 	_, err = stmt.Exec()
// 	if err != nil {
// 		log.Fatalf("exec failed: %s", err)
// 	}
// }
