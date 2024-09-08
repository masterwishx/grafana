import { MetricFindValue } from '@grafana/data';
import { AdHocFiltersVariable, CustomVariable, sceneGraph, SceneObject } from '@grafana/scenes';

import { VAR_OTEL_DEPLOYMENT_ENV, VAR_OTEL_RESOURCES } from '../shared';

import { OtelResourcesObject } from './types';

export const blessedList = (): Record<string, number> => {
  return {
    cloud_availability_zone: 0,
    cloud_region: 0,
    container_name: 0,
    k8s_cluster_name: 0,
    k8s_container_name: 0,
    k8s_cronjob_name: 0,
    k8s_daemonset_name: 0,
    k8s_deployment_name: 0,
    k8s_job_name: 0,
    k8s_namespace_name: 0,
    k8s_pod_name: 0,
    k8s_replicaset_name: 0,
    k8s_statefulset_name: 0,
    service_instance_id: 0,
    service_name: 0,
    service_namespace: 0,
  };
};

export function sortResources(resources: MetricFindValue[], excluded: string[]) {
  // these may be filtered
  const promotedList = blessedList();

  const blessed = Object.keys(promotedList);

  resources = resources.filter((resource) => {
    // if not in the list keep it
    const val = (resource.value ?? '').toString();

    if (!blessed.includes(val)) {
      return true;
    }
    // remove blessed filters
    // but indicate which are available
    promotedList[val] = 1;
    return false;
  });

  const promotedResources = Object.keys(promotedList)
    .filter((resource) => promotedList[resource] && !excluded.includes(resource))
    .map((v) => ({ text: v }));

  // put the filters first
  return promotedResources.concat(resources);
}

/**
 * Return a collection of labels and labels filters.
 * This data is used to build the join query to filter with otel resources
 *
 * @param otelResourcesObject
 * @returns a string that is used to add a join query to filter otel resources
 */
export function getOtelJoinQuery(otelResourcesObject: OtelResourcesObject): string {
  let otelResourcesJoinQuery = '';
  if (otelResourcesObject.filters && otelResourcesObject.labels) {
    // add support for otel data sources that are not standardized, i.e., have non unique target_info series by job, instance
    otelResourcesJoinQuery = `* on (job, instance) group_left(${otelResourcesObject.labels}) topk by (job, instance) (1, target_info{${otelResourcesObject.filters}})`;
  }

  return otelResourcesJoinQuery;
}

/**
 * Returns an object containing all the filters for otel resources as well as a list of labels
 *
 * @param scene
 * @param firstQueryVal
 * @returns
 */
export function getOtelResourcesObject(scene: SceneObject, firstQueryVal?: string): OtelResourcesObject {
  const otelResources = sceneGraph.lookupVariable(VAR_OTEL_RESOURCES, scene);
  // add deployment env to otel resource filters
  const otelDepEnv = sceneGraph.lookupVariable(VAR_OTEL_DEPLOYMENT_ENV, scene);

  let otelResourcesObject = { labels: '', filters: '' };

  if (otelResources instanceof AdHocFiltersVariable && otelDepEnv instanceof CustomVariable) {
    // get the collection of adhoc filters
    const otelFilters = otelResources.state.filters;

    // get the value for deployment_environment variable
    let otelDepEnvValue = String(otelDepEnv.getValue());
    // check if there are multiple environments
    const isMulti = otelDepEnvValue.includes(',');
    // start with the default label filters for deployment_environment
    let op = '=';
    let val = firstQueryVal ? firstQueryVal : otelDepEnvValue;
    // update the filters if multiple deployment environments selected
    if (isMulti) {
      op = '=~';
      val = val.split(',').join('|');
    }

    // start with the deployment environment
    let allFilters = `deployment_environment${op}"${val}"`;
    let allLabels = 'deployment_environment';

    // add the other OTEL resource filters
    for (let i = 0; i < otelFilters?.length; i++) {
      const labelName = otelFilters[i].key;
      const op = otelFilters[i].operator;
      const labelValue = otelFilters[i].value;

      allFilters += `,${labelName}${op}"${labelValue}"`;

      const addLabelToGroupLeft = labelName !== 'job' && labelName !== 'instance';

      if (addLabelToGroupLeft) {
        allLabels += `,${labelName}`;
      }
    }

    otelResourcesObject.labels = allLabels;
    otelResourcesObject.filters = allFilters;

    return otelResourcesObject;
  }
  return otelResourcesObject;
}