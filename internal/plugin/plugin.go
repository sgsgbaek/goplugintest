package plugin

import (
	"fmt"
	"os"

	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	"github.com/argoproj/argo-rollouts/rollout/trafficrouting/plugin/rpc"
	pluginTypes "github.com/argoproj/argo-rollouts/utils/plugin/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/globalaccelerator"
	"github.com/sirupsen/logrus"
)

var _ rpc.TrafficRouterPlugin = &RpcPlugin{}

type RpcPlugin struct {
	LogCtx *logrus.Entry
}

type Environment struct {
	// acceleratorArn   string
	endpointGroupArn string
}

var env = new(Environment)

func (p *Environment) GetEnvs() {
	// p.acceleratorArn = os.Getenv("ACCELERATOR_ARN")
	p.endpointGroupArn = os.Getenv("ENDPOINT_GROUP_ARN")
}

func (p *RpcPlugin) InitPlugin() pluginTypes.RpcError {
	p.LogCtx.Info("InitPlugin-SGSGBAEK")
	env.GetEnvs()
	// p.LogCtx.Debug("acceleratorArn : %s", env.acceleratorArn)
	p.LogCtx.Debug("endpointGroupArn : %s", env.endpointGroupArn)
	return pluginTypes.RpcError{}
}

func updateListenerTraffic(endpointGroupArn string, trafficPercentage float64) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))
	svc := globalaccelerator.New(sess)
	input := &globalaccelerator.UpdateEndpointGroupInput{
		EndpointGroupArn:      aws.String(endpointGroupArn),
		TrafficDialPercentage: aws.Float64(trafficPercentage),
	}
	_, err := svc.UpdateEndpointGroup(input)
	if err != nil {
		return fmt.Errorf("failed to update endpoint group: %v", err)
	}
	return nil
}

// SetWeight modifies Nginx Ingress resources to reach desired state
func (r *RpcPlugin) SetWeight(ro *v1alpha1.Rollout, desiredWeight int32, additionalDestinations []v1alpha1.WeightDestination) pluginTypes.RpcError {
	r.LogCtx.Info("PRC set weight called with weight: %d", desiredWeight)
	updateListenerTraffic(env.endpointGroupArn, float64(desiredWeight))
	return pluginTypes.RpcError{}
}

func (r *RpcPlugin) SetHeaderRoute(ro *v1alpha1.Rollout, headerRouting *v1alpha1.SetHeaderRoute) pluginTypes.RpcError {
	return pluginTypes.RpcError{}
}

func (r *RpcPlugin) VerifyWeight(ro *v1alpha1.Rollout, desiredWeight int32, additionalDestinations []v1alpha1.WeightDestination) (pluginTypes.RpcVerified, pluginTypes.RpcError) {
	return pluginTypes.NotImplemented, pluginTypes.RpcError{}
}

// UpdateHash informs a traffic routing reconciler about new canary/stable pod hashes
func (r *RpcPlugin) UpdateHash(ro *v1alpha1.Rollout, canaryHash, stableHash string, additionalDestinations []v1alpha1.WeightDestination) pluginTypes.RpcError {
	return pluginTypes.RpcError{}
}

func (r *RpcPlugin) SetMirrorRoute(ro *v1alpha1.Rollout, setMirrorRoute *v1alpha1.SetMirrorRoute) pluginTypes.RpcError {
	return pluginTypes.RpcError{}
}

func (r *RpcPlugin) RemoveManagedRoutes(ro *v1alpha1.Rollout) pluginTypes.RpcError {
	return pluginTypes.RpcError{}
}

func (r *RpcPlugin) Type() string {
	return "plugin-globalaccelerator"
}

func DoAction() string {
	return "Plugin action executed - sgsgbaek\n"
}
