package loki

type Client interface {
	Debugf(job, template string, args ...interface{})
	Infof(job, template string, args ...interface{})
	Warnf(job, template string, args ...interface{})
	Errorf(job, template string, args ...interface{})
}
