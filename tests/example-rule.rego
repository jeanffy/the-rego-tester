package example.opa

default allowed := false
allowed if {
  print("allowed", input.userName)
  input.userName == "alice"
}

state := "inactive" if {
  print("state 1", input.numberOfVisits)
  input.numberOfVisits <= 5
}

state := "active" if {
  print("state 2", input.numberOfVisits)
  input.numberOfVisits > 5
  input.numberOfVisits < 100
}

state := "super-active" if {
  print("state 3", input.numberOfVisits)
  input.numberOfVisits > 100
}
