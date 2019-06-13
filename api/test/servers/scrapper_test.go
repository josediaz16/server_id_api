package servers_test

import (
  "testing"
  "server_id_api/servers"
)

func TestGetDomainHeadIsOk(t *testing.T) {
  title, logo := servers.GetDomainHead("stackoverflow.com")

  expectedTitle := "Stack Overflow - Where Developers Learn, Share, & Build Careers"
  expectedLogo := "https://cdn.sstatic.net/Sites/stackoverflow/img/favicon.ico?v=4f32ecc8f43d"

  if title != expectedTitle {
    t.Errorf("TestGetDomainHeadIsOk(stackoverflow.com) got title %v, should be %v", title, expectedTitle)
  }

  if logo != expectedLogo {
    t.Errorf("TestGetDomainHeadIsOk(stackoverflow.com) got logo %v, should be %v", logo, expectedLogo)
  }
}

func TestCheckRelativePathOnlyPath(t *testing.T) {
  logo := servers.CheckRelativePath("instagram.com", "/assets/favicon-b8.ico")

  expectedLogo := "http://instagram.com/assets/favicon-b8.ico"

  if logo != expectedLogo {
    t.Errorf("TestCheckRelativePath got logo %v, should be %v", logo, expectedLogo)
  }
}

func TestCheckRelativePathFullUrl(t *testing.T) {
  logo := servers.CheckRelativePath("instagram.com", "http://instagram.com/assets/favicon-b8.ico")

  expectedLogo := "http://instagram.com/assets/favicon-b8.ico"

  if logo != expectedLogo {
    t.Errorf("TestCheckRelativePath got logo %v, should be %v", logo, expectedLogo)
  }
}
