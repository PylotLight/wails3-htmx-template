/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["../components/*.{go,templ}", "../*.{go,templ}", './pages/**/*.{html,js}', "./index.html", "./public/systray/*.html", './src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}',],
    theme: {
        extend: {
            // fontFamily: {
            //     sans: ['Inter'], // Add remaining fonts from :root
            // },
            fontSize: {
                base: '16px', // Set base font size
            },
            lineHeight: {
                base: '24px', // Set base line height
            },
            fontWeight: {
                normal: '400', // Set normal font weight
            },
            textColor: {
                primary: 'rgba(255, 255, 255, 0.87)', // Primary text color
                'dark': 'rgba(255, 255, 255, 0.87)', // Dark mode text color (if applicable)
            },
            backgroundColor: {
                body: 'rgba(27, 38, 54, 1)', // Background color
            },
            // ... other theme extensions
        },
    },
    plugins: [require("daisyui")],
}