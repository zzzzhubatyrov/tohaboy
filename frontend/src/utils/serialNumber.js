/**
 * Generates a serial number based on the equipment category and current date
 * @param {string} categoryName - Name of the equipment category
 * @returns {string} Generated serial number
 */
export function generateSerialNumber(categoryName) {
    if (!categoryName) return ''

    // Get current date components
    const date = new Date()
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')

    // Get category prefix (first 3 letters in uppercase)
    const prefix = categoryName
        .replace(/[^a-zA-Zа-яА-Я]/g, '') // Remove non-letters
        .slice(0, 3)
        .toUpperCase()

    // Generate random number (3 digits)
    const random = Math.floor(Math.random() * 1000).toString().padStart(3, '0')

    // Combine all parts
    return `${prefix}-${year}${month}-${random}`
} 