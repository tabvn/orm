package orm

import "log"

func (db *DB) Models(models ...interface{}) {
	if models != nil {
		for _, m := range models {
			if model, err := NewModel(m); err != nil {
				log.Fatal(err)
				return
			} else {
				db.models = append(db.models, model)
			}
		}
	}
}
