// go2hal v0.3.0
// Copyright (c) 2017 Patrick Moule
// License: MIT

package hal

// ResourceFactory is a helper for creating resources and links.
type ResourceFactory interface {
	CreateRootResource(href string, method string) Resource
	CreateEmbeddedResource(href string, method string) Resource
	CreateLink(relationName string, href string, method string, curieLinkName string) LinkRelation
	CreateResourceLink(relationName string, curieLinkName string) ResourceRelation
}

type resourceFactory struct {
	curieLinks map[string]*LinkObject
}

// NewResourceFactory initialises a ResourceFactory with a set of CURIE links.
func NewResourceFactory(curieLinks []*LinkObject) ResourceFactory {
	links := make(map[string]*LinkObject)

	for _, link := range curieLinks {
		links[link.Name] = link
	}

	factory := &resourceFactory{curieLinks: links}
	return factory
}

// CreateLink creates a Link Relation with given relationName and href. A CURIE link can
// be added by curieLinkName. The real CURIE link is picked from the set of CURIE links the factory
// is initialised with.
func (rf *resourceFactory) CreateLink(relationName string, href string, method string, curieLinkName string) LinkRelation {
	link, linkError := NewLinkObject(href, method)

	if linkError != nil {
		return nil
	}

	relation, relationError := NewLinkRelation(relationName)

	if relationError != nil {
		return nil
	}

	relation.SetLink(link)

	if curieLinkName != "" {
		curieLink := rf.curieLinks[curieLinkName]

		if curieLink != nil {
			relation.SetCurieLink(curieLink)
		}
	}

	return relation
}

// CreateResourceLink creates a Link Relation with given relationName. A CURIE link can
// be added by curieLinkName. The real CURIE link is picked from the set of CURIE links the factory
// is initialised with.
func (rf *resourceFactory) CreateResourceLink(relationName string, curieLinkName string) ResourceRelation {
	relation, relationError := NewResourceRelation(relationName)

	if relationError != nil {
		return nil
	}

	if curieLinkName != "" {
		curieLink := rf.curieLinks[curieLinkName]

		if curieLink != nil {
			relation.SetCurieLink(curieLink)
		}
	}

	return relation
}

// CreateRootResource creates a root Resource with self link from given href.
// Additionally all CURIE links given at ResourceFactory initialisation are added.
func (rf *resourceFactory) CreateRootResource(href string, method string) Resource {
	resource := rf.createResource(href, method)

	curieLinks := []*LinkObject{}

	for _, v := range rf.curieLinks {
		curieLinks = append(curieLinks, v)
	}

	resource.AddCurieLinks(curieLinks)

	return resource
}

// CreateEmbeddedResource creates an embedded Resource with self link from given href.
func (rf *resourceFactory) CreateEmbeddedResource(href string, method string) Resource {
	resource := rf.createResource(href, method)

	return resource
}

func (rf *resourceFactory) createResource(href string, method string) Resource {
	resource := NewResourceObject()
	selfLink, err := NewLinkObject(href, method)

	if err != nil {
		return resource
	}

	self := NewSelfLinkRelation()
	self.SetLink(selfLink)
	resource.AddLink(self)

	return resource
}
