import { gte, SemVer } from 'semver';

import { isMetricAggregationWithField } from './components/QueryEditor/MetricAggregationsEditor/aggregations';
import { metricAggregationConfig } from './components/QueryEditor/MetricAggregationsEditor/utils';
import { ElasticsearchQuery, MetricAggregation, MetricAggregationWithInlineScript } from './types';
import { DataFrame, FieldType, LogLevel } from '@grafana/data';

export const describeMetric = (metric: MetricAggregation) => {
  if (!isMetricAggregationWithField(metric)) {
    return metricAggregationConfig[metric.type].label;
  }

  // TODO: field might be undefined
  return `${metricAggregationConfig[metric.type].label} ${metric.field}`;
};

/**
 * Utility function to clean up aggregations settings objects.
 * It removes nullish values and empty strings, array and objects
 * recursing over nested objects (not arrays).
 * @param obj
 */
export const removeEmpty = <T extends {}>(obj: T): Partial<T> =>
  Object.entries(obj).reduce((acc, [key, value]) => {
    // Removing nullish values (null & undefined)
    if (value == null) {
      return { ...acc };
    }

    // Removing empty arrays (This won't recurse the array)
    if (Array.isArray(value) && value.length === 0) {
      return { ...acc };
    }

    // Removing empty strings
    if (typeof value === 'string' && value.length === 0) {
      return { ...acc };
    }

    // Recursing over nested objects
    if (!Array.isArray(value) && typeof value === 'object') {
      const cleanObj = removeEmpty(value);

      if (Object.keys(cleanObj).length === 0) {
        return { ...acc };
      }

      return { ...acc, [key]: cleanObj };
    }

    return {
      ...acc,
      [key]: value,
    };
  }, {});

/**
 *  This function converts an order by string to the correct metric id For example,
 *  if the user uses the standard deviation extended stat for the order by,
 *  the value would be "1[std_deviation]" and this would return "1"
 */
export const convertOrderByToMetricId = (orderBy: string): string | undefined => {
  const metricIdMatches = orderBy.match(/^(\d+)/);
  return metricIdMatches ? metricIdMatches[1] : void 0;
};

/** Gets the actual script value for metrics that support inline scripts.
 *
 *  This is needed because the `script` is a bit polymorphic.
 *  when creating a query with Grafana < 7.4 it was stored as:
 * ```json
 * {
 *    "settings": {
 *      "script": {
 *        "inline": "value"
 *      }
 *    }
 * }
 * ```
 *
 * while from 7.4 it's stored as
 * ```json
 * {
 *    "settings": {
 *      "script": "value"
 *    }
 * }
 * ```
 *
 * This allows us to access both formats and support both queries created before 7.4 and after.
 */
export const getScriptValue = (metric: MetricAggregationWithInlineScript) =>
  (typeof metric.settings?.script === 'object' ? metric.settings?.script?.inline : metric.settings?.script) || '';

export const isSupportedVersion = (version: SemVer): boolean => {
  if (gte(version, '7.16.0')) {
    return true;
  }

  return false;
};

export const unsupportedVersionMessage =
  'Support for Elasticsearch versions after their end-of-life (currently versions < 7.16) was removed. Using unsupported version of Elasticsearch may lead to unexpected and incorrect results.';

// To be considered a time series query, the last bucked aggregation must be a Date Histogram
export const isTimeSeriesQuery = (query: ElasticsearchQuery): boolean => {
  return query?.bucketAggs?.slice(-1)[0]?.type === 'date_histogram';
};

/*
 * This regex matches 3 types of variable reference with an optional format specifier
 * There are 6 capture groups that replace will return
 * \$(\w+)                                    $var1
 * \[\[(\w+?)(?::(\w+))?\]\]                  [[var2]] or [[var2:fmt2]]
 * \${(\w+)(?:\.([^:^\}]+))?(?::([^\}]+))?}   ${var3} or ${var3.fieldPath} or ${var3:fmt3} (or ${var3.fieldPath:fmt3} but that is not a separate capture group)
 */
export const variableRegex = /\$(\w+)|\[\[(\w+?)(?::(\w+))?\]\]|\${(\w+)(?:\.([^:^\}]+))?(?::([^\}]+))?}/g;

export const logLevelMap: Record<string, LogLevel> = {
  emerg: LogLevel.emerg,
  fatal: LogLevel.fatal,
  alert: LogLevel.alert,
  crit: LogLevel.crit,
  critical: LogLevel.critical,
  warn: LogLevel.warn,
  warning: LogLevel.warning,
  err: LogLevel.err,
  eror: LogLevel.eror,
  error: LogLevel.error,
  info: LogLevel.info,
  information: LogLevel.information,
  informational: LogLevel.informational,
  notice: LogLevel.notice,
  dbug: LogLevel.dbug,
  debug: LogLevel.debug,
  trace: LogLevel.trace,
  unknown: LogLevel.unknown,
}

export function dataFrameLogLevel(dataFrame: DataFrame): LogLevel {
  const field = dataFrame.fields.find((f) => f.type === FieldType.number);
  const level = field?.labels?.['level'] ?? field?.labels?.['lvl'] ?? field?.labels?.['loglevel'] ?? '';
  if (!level) {
    return LogLevel.unknown;
  }
  return level in logLevelMap ? logLevelMap[level] : LogLevel.unknown;
}
