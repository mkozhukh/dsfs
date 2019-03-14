package store

import (
	"encoding/json"
	"io/ioutil"
)

func NewUserCollection(path string) (*UserCollection, error) {
	data := UserCollection{path: path}
	err := data.loadFromFile()

	return &data, err
}

type UserCollection struct {
	path  string
	pull  map[uint]*User
	order []User

	counter uint
}

type userFormat struct {
	Order   []User
	Counter uint
}

func (c *UserCollection) Get(id uint) *User {
	item, ok := c.pull[id]
	if !ok {
		return &User{}
	}

	return item
}

func (c *UserCollection) Exists(id uint) bool {
	_, ok := c.pull[id]
	return ok
}

func (c *UserCollection) GetAll() []User {
	return c.order
}

type UserLocator func(*User) bool

func (c *UserCollection) First(loc UserLocator) *User {
	for i := range c.order {
		if loc(&c.order[i]) {
			return &c.order[i]
		}
	}

	return &User{}
}

func (c *UserCollection) Save(obj *User) error {
	index := -1
	if obj.ID == 0 {
		c.counter++
		obj.ID = c.counter
	} else {
		index = c.indexOf(obj.ID)
	}

	if index == -1 {
		c.order = append(c.order, *obj)
	} else {
		c.order[index] = *obj
	}

	c.pull[obj.ID] = obj

	return c.saveToFile()
}

func (c *UserCollection) Delete(id uint) error {
	index := c.indexOf(id)
	if index == -1 {
		return nil
	}

	c.order = append(c.order[:index], c.order[index+1:]...)
	delete(c.pull, id)

	return c.saveToFile()
}

func (c *UserCollection) indexOf(key uint) int {
	for i := range c.order {
		if c.order[i].ID == key {
			return i
		}
	}

	return -1
}

func (c *UserCollection) loadFromFile() error {
	c.pull = make(map[uint]*User)

	bytes, err := ioutil.ReadFile(c.path)
	if err != nil {
		return err
	}

	temp := userFormat{}
	err = json.Unmarshal(bytes, &temp)
	if err != nil {
		return err
	}

	c.order = temp.Order
	c.counter = temp.Counter
	for i := range c.order {
		c.pull[c.order[i].ID] = &c.order[i]
	}

	return nil
}

func (c *UserCollection) saveToFile() error {
	bytes, err := json.Marshal(userFormat{
		Order:   c.order,
		Counter: c.counter,
	})

	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.path, bytes, 0644)
}
