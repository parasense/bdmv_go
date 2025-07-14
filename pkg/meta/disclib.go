package meta

/*
	Remarks:

	There is something called DiscLib.
	It's an XML scheme for categorizing BluRay metadata.
	The schema is probably not available publicly, or otherwise obscure.
	Therefore everything here is a reverse enginering effort.

*/

import (
	"encoding/xml"
)

// DiscLib represents the root <disclib> element
type DiscLib struct {
	XMLName  xml.Name `xml:"urn:BDA:bdmv;disclib disclib"`
	DiscInfo DiscInfo `xml:"urn:BDA:bdmv;discinfo discinfo"`
}

// DiscInfo represents the <di:discinfo> element
type DiscInfo struct {
	XMLName     xml.Name    `xml:"urn:BDA:bdmv;discinfo discinfo"`
	Title       Title       `xml:"title"`
	Description Description `xml:"description"`
	Language    *string     `xml:"language"` // Optional field
}

// Title represents the <di:title> element
type Title struct {
	Name      string `xml:"name"`
	NumSets   *int   `xml:"numSets"`   // Optional field
	SetNumber *int   `xml:"setNumber"` // Optional field
}

// Description represents the <di:description> element
type Description struct {
	Thumbnails      []Thumbnail      `xml:"thumbnail"`
	TableOfContents *TableOfContents `xml:"tableOfContents"` // Optional field
}

// Thumbnail represents the <di:thumbnail> element
type Thumbnail struct {
	Href string  `xml:"href,attr"`
	Size *string `xml:"size,attr"` // Optional attribute
}

// TableOfContents represents the <di:tableOfContents> element
type TableOfContents struct {
	TitleNames []TitleName `xml:"titleName"`
}

// TitleName represents the <di:titleName> element
type TitleName struct {
	TitleNumber string `xml:"titleNumber,attr"`
	Name        string `xml:",chardata"`
}
