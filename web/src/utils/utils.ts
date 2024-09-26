// ------------------------------------------------------------
// : Import
// ------------------------------------------------------------

import _ from 'lodash'
// ------------------------------------------------------------
// : Alias
// ------------------------------------------------------------
export const once       = _.once
export const every      = _.every

export const get        = _.get
export const set        = _.set

export const map        = _.map
export const each       = _.each
export const find       = _.find
export const filter     = _.filter
export const is_empty   = _.isEmpty
export const sort_by    = _.sortBy
export const find_index = _.findIndex

export const trim       = _.trim
export const is_array   = _.isArray

// ------------------------------------------------------------
// : Functions
// ------------------------------------------------------------
export function is_debug() {
    return import.meta.env.MODE === 'development'
}
