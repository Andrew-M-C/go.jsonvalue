package exporter

import (
	"reflect"
	"sync"
)

// g stands for Global
type g struct {
	lock sync.RWMutex

	exportersByType map[reflect.Type]Exporter

	debugf func(string, ...any)
}

var internal = &g{}

func init() {
	internal.exportersByType = make(map[reflect.Type]Exporter)
	internal.debugf = func(string, ...any) {}
}

func (i *g) loadExportersByType(typ reflect.Type) (Exporter, bool) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	e, exist := i.exportersByType[typ]
	return e, exist
}

func (i *g) storeExporterByType(typ reflect.Type, e Exporter) {
	internal.lock.Lock()
	defer internal.lock.Unlock()

	internal.exportersByType[typ] = e
}
