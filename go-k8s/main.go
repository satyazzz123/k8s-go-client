// package main

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"net/http"
// 	"net/http/httputil"
// 	"os"

// 	"k8s.io/apimachinery/pkg/api/errors"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/runtime/schema"
// 	"k8s.io/client-go/dynamic"
// 	"k8s.io/client-go/tools/clientcmd"
// )

// // Tracer implements http.RoundTripper.  It prints each request and
// // response/error to os.Stderr.  WARNING: this may output sensitive information
// // including bearer tokens.
// type Tracer struct {
// 	http.RoundTripper
// }

// // RoundTrip calls the nested RoundTripper while printing each request and
// // response/error to os.Stderr on either side of the nested call.  WARNING: this
// // may output sensitive information including bearer tokens.
// func (t *Tracer) RoundTrip(req *http.Request) (*http.Response, error) {
// 	// Dump the request to os.Stderr.
// 	b, err := httputil.DumpRequestOut(req, true)
// 	if err != nil {
// 		return nil, err
// 	}
// 	os.Stderr.Write(b)
// 	os.Stderr.Write([]byte{'\n'})

// 	// Call the nested RoundTripper.
// 	resp, err := t.RoundTripper.RoundTrip(req)

// 	// If an error was returned, dump it to os.Stderr.
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		return resp, err
// 	}

// 	// Dump the response to os.Stderr.
// 	b, err = httputil.DumpResponse(resp, req.URL.Query().Get("watch") != "true")
// 	if err != nil {
// 		return nil, err
// 	}
// 	os.Stderr.Write(b)
// 	os.Stderr.Write([]byte{'\n'})

// 	return resp, err
// }

// // go run main.go -namespace kube-system
// func main() {
// 	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
// 	if len(defaultKubeconfig) == 0 {
// 		defaultKubeconfig = clientcmd.RecommendedHomeFile
// 	}

// 	kubeconfig := flag.String(clientcmd.RecommendedConfigPathFlag,
// 		defaultKubeconfig, "absolute path to the kubeconfig file")

// 	namespace := flag.String("namespace", metav1.NamespaceDefault,
// 		"create the deployment in this namespace")

// 	verbose := flag.Bool("verbose", false, "display HTTP calls")

// 	flag.Parse()

// 	rc, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	if *verbose {
// 		// wrap the default RoundTripper with an instance of Tracer.
// 		rc.WrapTransport = func(rt http.RoundTripper) http.RoundTripper {
// 			return &Tracer{rt}
// 		}
// 	}

// 	// create a new dynamic client using the rest.Config
// 	dc, err := dynamic.NewForConfig(rc)
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	// identify pods resource
// 	// gvr := schema.GroupVersionResource{
// 	// 	Version:  "v1",
// 	// 	Resource: "pods",
// 	// }
// 	gvr := schema.GroupVersionResource{
// 		Group:    "apps",
// 		Version:  "v1",
// 		Resource: "deployments",
// 	}

// 	// list all pods in the specified namespace
// 	res, err := dc.Resource(gvr).
// 		Namespace(*namespace).
// 		List(context.TODO(), metav1.ListOptions{})
// 	if err != nil {
// 		if !errors.IsNotFound(err) {
// 			panic(err)
// 		}
// 	}

// 	// for each pod, print just the name
// 	for _, el := range res.Items {
// 		fmt.Printf("%v\n", el.GetName())
// 	}
// }

// package main

// import (
// 	"fmt" // Importing the fmt package for printing messages

// 	"k8s.io/apimachinery/pkg/labels"    // Importing the labels package from Kubernetes
// 	"k8s.io/apimachinery/pkg/selection" // Importing the selection package from Kubernetes
// )

// const (
// 	keySupportBy = "app.kubernetes.io/support-by" // Define a constant for the label key
// 	keyPartOf    = "app.kubernetes.io/part-of"    // Define a constant for another label key
// )

// func main() {
// 	// Create a requirement for the existence of the 'app.kubernetes.io/support-by' label
// 	supportExists, err := labels.NewRequirement(keySupportBy, selection.Exists, []string{})
// 	if err != nil {
// 		panic(err) // If there's an error, panic (stop execution)
// 	}

// 	// Create a requirement where the 'app.kubernetes.io/support-by' label value must be either 'team_1' or 'team_2'
// 	supportTeam, err := labels.NewRequirement(keySupportBy, selection.In, []string{"team_1", "team_2"})
// 	if err != nil {
// 		panic(err) // If there's an error, panic (stop execution)
// 	}

// 	// Create a requirement where the 'app.kubernetes.io/part-of' label value must be 'payment-system'
// 	partOfReq, err := labels.NewRequirement(keyPartOf, selection.Equals, []string{"payment-system"})
// 	if err != nil {
// 		panic(err) // If there's an error, panic (stop execution)
// 	}

// 	// Create an empty label selector
// 	selector := labels.NewSelector()

// 	// Add the requirements to the selector
// 	selector = selector.Add(*supportExists, *supportTeam, *partOfReq)

// 	// Print the selector (it will be rendered as in `kubectl label...`)
// 	fmt.Printf("here your selector: %s\n", selector.String())
// }

// package main

// import (
// 	"context" // Import the context package for managing context
// 	"flag"    // Import the flag package for command-line flag parsing
// 	"fmt"     // Import the fmt package for printing messages
// 	"os"      // Import the os package for environment variables
// 	"strings" // Import the strings package for string manipulation

// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // Import Kubernetes metav1 package for metav1.ListOptions
// 	"k8s.io/apimachinery/pkg/labels"              // Import Kubernetes labels package for label parsing
// 	"k8s.io/apimachinery/pkg/runtime/schema"      // Import Kubernetes schema package for GroupVersionResource
// 	"k8s.io/client-go/dynamic"                    // Import Kubernetes dynamic package for dynamic client
// 	"k8s.io/client-go/tools/clientcmd"            // Import Kubernetes clientcmd package for building client config
// )

// func main() {
// 	// Get the default kubeconfig file path
// 	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
// 	if len(defaultKubeconfig) == 0 {
// 		defaultKubeconfig = clientcmd.RecommendedHomeFile
// 	}

// 	// Define and parse command-line flags for the kubeconfig file path
// 	kubeconfig := flag.String(clientcmd.RecommendedConfigPathFlag,
// 		defaultKubeconfig, "absolute path to the kubeconfig file")
// 	flag.Parse()

// 	// Build a Kubernetes client configuration from the specified kubeconfig file path
// 	rc, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Create a dynamic client for accessing Kubernetes resources
// 	dc, err := dynamic.NewForConfig(rc)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Parse any additional arguments as a label selector and convert it into a string
// 	var filter string
// 	if len(flag.Args()) > 0 {
// 		sel, err := labels.Parse(strings.Join(flag.Args(), " "))
// 		if err != nil {
// 			panic(err)
// 		}
// 		filter = sel.String()
// 	}

// 	// Define the resource type to list (in this case, namespaces)
// 	gvr := schema.GroupVersionResource{
// 		Version:  "v1",
// 		Resource: "namespaces",
// 	}

// 	// Use the dynamic client to list namespaces with the specified label selector
// 	res, err := dc.Resource(gvr).
// 		List(context.TODO(), metav1.ListOptions{LabelSelector: filter})
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Print the names of the namespaces that match the label selector
// 	for _, el := range res.Items {
// 		fmt.Println(el.GetName())
// 	}
// }

package main

import (
	"encoding/json" // Import the encoding/json package for JSON encoding and decoding
	"fmt"           // Import the fmt package for printing messages

	"k8s.io/apimachinery/pkg/util/errors" // Import the errors package from Kubernetes
	"k8s.io/client-go/discovery"          // Import the discovery package from Kubernetes
	"k8s.io/client-go/tools/clientcmd"    // Import the clientcmd package from Kubernetes
)

func main() {
	// Create a client configuration loader using the default loading rules and overrides
	configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	// Get the client configuration from the loader
	rc, err := configLoader.ClientConfig()
	if err != nil {
		panic(err) // If there's an error, panic (stop execution)
	}

	// Create a new DiscoveryClient using the client configuration
	dc, err := discovery.NewDiscoveryClientForConfig(rc)
	if err != nil {
		panic(err) // If there's an error, panic (stop execution)
	}

	// Storage for errors
	errs := []error{}

	// Retrieve the supported resources with the version preferred by the server
	lists, err := dc.ServerPreferredResources()
	if err != nil {
		errs = append(errs, err)
	}

	// Utility struct holding information to print
	type info struct {
		Kind       string   `json:"kind"`
		APIVersion string   `json:"apiVersion"`
		Name       string   `json:"name"`
		Verbs      []string `json:"verbs"`
	}

	// Iterate all the APIResource collections
	for _, list := range lists {
		if len(list.APIResources) == 0 {
			continue
		}

		// Grab the API resource info
		for _, el := range list.APIResources {
			if len(el.Verbs) == 0 {
				continue
			}

			// Create a temporary struct with resource information
			tmp := info{el.Kind, list.GroupVersion, el.Name, el.Verbs}
			// Convert the struct to JSON
			res, err := json.Marshal(&tmp)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			// Print the JSON
			fmt.Printf("%s\n", res)
		}
	}

	// If there has been an error, print it on the screen
	if len(errs) > 0 {
		panic(errors.NewAggregate(errs))
	}
}
