export function isElectron(): boolean {
  return (
    typeof window !== 'undefined' &&
    typeof window.process === 'object' &&
    window.process.type === 'renderer'
  )
}
