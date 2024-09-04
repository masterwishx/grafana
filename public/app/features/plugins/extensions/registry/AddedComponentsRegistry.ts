import { PluginExtensionAddedComponentConfig } from '@grafana/data';

import { wrapWithPluginContext } from '../utils';
import {
  extensionPointEndsWithVersion,
  isExtensionPointIdValid,
  isGrafanaCoreExtensionPoint,
  isReactComponent,
} from '../validators';

import { PluginExtensionConfigs, Registry, RegistryType } from './Registry';

export type AddedComponentRegistryItem<Props = {}> = {
  pluginId: string;
  title: string;
  description: string;
  component: React.ComponentType<Props>;
};

export class AddedComponentsRegistry extends Registry<
  AddedComponentRegistryItem[],
  PluginExtensionAddedComponentConfig
> {
  constructor(initialState: RegistryType<AddedComponentRegistryItem[]> = {}) {
    super({
      initialState,
    });
  }

  mapToRegistry(
    registry: RegistryType<AddedComponentRegistryItem[]>,
    item: PluginExtensionConfigs<PluginExtensionAddedComponentConfig>
  ): RegistryType<AddedComponentRegistryItem[]> {
    const { pluginId, configs } = item;

    for (const config of configs) {
      const log = this.logger.child({
        description: config.description,
        title: config.title,
        pluginId,
      });

      if (!isReactComponent(config.component)) {
        log.warning(
          `Could not register added component with title '${config.title}'. Reason: The provided component is not a valid React component.`
        );
        continue;
      }

      if (!config.title) {
        log.warning(`Could not register added component with title '${config.title}'. Reason: Title is missing.`);
        continue;
      }

      if (!config.description) {
        log.warning(`Could not register added component with title '${config.title}'. Reason: Description is missing.`);
        continue;
      }

      const extensionPointIds = Array.isArray(config.targets) ? config.targets : [config.targets];
      for (const extensionPointId of extensionPointIds) {
        const pointIdLog = log.child({ id: extensionPointId });

        if (!isExtensionPointIdValid(pluginId, extensionPointId)) {
          pointIdLog.warning(
            `Could not register added component with id '${extensionPointId}'. Reason: The component id does not match the id naming convention. Id should be prefixed with plugin id or grafana. e.g '<grafana|myorg-basic-app>/my-component-id/v1'.`
          );
          continue;
        }

        if (!isGrafanaCoreExtensionPoint(extensionPointId) && !extensionPointEndsWithVersion(extensionPointId)) {
          pointIdLog.warning(
            `Added component with id '${extensionPointId}' does not match the convention. It's recommended to suffix the id with the component version. e.g 'myorg-basic-app/my-component-id/v1'.`
          );
        }

        const result = {
          pluginId,
          component: wrapWithPluginContext(pluginId, config.component),
          description: config.description,
          title: config.title,
        };

        if (!(extensionPointId in registry)) {
          registry[extensionPointId] = [result];
        } else {
          registry[extensionPointId].push(result);
        }
      }
    }

    return registry;
  }
}