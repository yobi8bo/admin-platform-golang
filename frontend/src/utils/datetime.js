export function formatDateTime(value) {
  if (!value) return ''
  if (typeof value === 'string' && Number.isNaN(Number(value))) return value

  const date = new Date(Number(value))
  if (Number.isNaN(date.getTime())) return value

  const pad = (num) => String(num).padStart(2, '0')
  return [
    date.getFullYear(),
    pad(date.getMonth() + 1),
    pad(date.getDate())
  ].join('-') + ` ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}
