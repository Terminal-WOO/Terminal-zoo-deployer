export const cpuPattern = /^(\d+m|\d+(\.\d+)?)$/
export const memoryPattern = /^\d+(Mi|Gi)$/

export function validateResources(r: { cpuRequests:string; cpuLimits:string; memoryRequests:string; memoryLimits:string }): boolean {
  return (
    cpuPattern.test(r.cpuRequests) &&
    cpuPattern.test(r.cpuLimits) &&
    memoryPattern.test(r.memoryRequests) &&
    memoryPattern.test(r.memoryLimits)
  )
}