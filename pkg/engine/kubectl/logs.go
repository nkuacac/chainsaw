package kubectl

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmespath-community/go-jmespath/pkg/binding"
	"github.com/kyverno/chainsaw/pkg/apis/v1alpha1"
)

func Logs(ctx context.Context, tc binding.Bindings, collector *v1alpha1.PodLogs) (string, []string, error) {
	if collector == nil {
		return "", nil, errors.New("collector is null")
	}
	name, err := collector.Name.Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	namespace, err := collector.Namespace.Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	selector, err := collector.Selector.Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	container, err := collector.Container.Value(ctx, tc)
	if err != nil {
		return "", nil, err
	}
	if name == "" && selector == "" {
		return "", nil, errors.New("a name or selector must be specified")
	}
	if name != "" && selector != "" {
		return "", nil, errors.New("name cannot be provided when a selector is specified")
	}
	args := []string{"logs", "--prefix"}
	if name != "" {
		args = append(args, name)
	} else if selector != "" {
		args = append(args, "-l", selector)
	}
	if namespace == "" {
		namespace = "$NAMESPACE"
	}
	args = append(args, "-n", namespace)
	if container == "" {
		args = append(args, "--all-containers")
	} else {
		args = append(args, "-c", container)
	}
	if collector.Tail != nil {
		args = append(args, "--tail", fmt.Sprint(*collector.Tail))
	}
	return "kubectl", args, nil
}
