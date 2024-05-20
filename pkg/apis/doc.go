package eksamazonawscom

import (
	_ "embed"

	"github.com/awslabs/operatorpkg/object"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

//go:generate controller-gen crd object paths="./..." output:crd:artifacts:config=crds/templates
var (
	//go:embed crds/templates/eks.amazonaws.com_cninodes.yaml
	CRDBytes []byte

	CRDs = []*apiextensionsv1.CustomResourceDefinition{
		object.Unmarshal[apiextensionsv1.CustomResourceDefinition](CRDBytes),
	}
)
