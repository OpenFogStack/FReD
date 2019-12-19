package app

// Keygroups is an interface that abstracts the component that stores Keygroups metadata.
type Keygroups interface {
	Create(kg string) error
	Delete(kg string) error
	Exists(kg string) bool
}

// createKeygroup creates a new keygroup with the specified name in Storage.
func (a *App) createKeygroup(kgname string) error {
	err := a.kg.Create(kgname)

	if err != nil {
		return err
	}

	err = a.sd.CreateKeygroup(kgname)

	if err != nil {
		return err
	}

	return nil
}

// deleteKeygroup removes the keygroup with the specified name from Storage.
func (a *App) deleteKeygroup(kgname string) error {
	err := a.kg.Delete(kgname)

	if err != nil {
		return err
	}

	err = a.sd.DeleteKeygroup(kgname)

	if err != nil {
		return err
	}

	return nil
}
