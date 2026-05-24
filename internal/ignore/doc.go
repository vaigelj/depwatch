// Package ignore manages a list of dependency suppressions that allow
// users to acknowledge known issues and prevent them from appearing in
// future depwatch reports.
//
// Suppressions are loaded from a JSON file (typically .depignore at the
// repository root) and can carry an optional expiry date so that
// temporary waivers are automatically re-activated once they lapse.
//
// Example .depignore file:
//
//	{
//	  "ignore": [
//	    {
//	      "ecosystem": "npm",
//	      "package":   "lodash",
//	      "reason":    "internal fork, patched separately",
//	      "expires":   "2025-12-31T00:00:00Z"
//	    }
//	  ]
//	}
package ignore
