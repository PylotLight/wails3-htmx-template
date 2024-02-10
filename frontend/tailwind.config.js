/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["../components/*.{go,templ}", "../*.{go,templ}", './pages/**/*.{html,js}', "./index.html", "./public/systray/*.html", './src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}',],
    theme: {
        extend: {},
    },
    plugins: [require("daisyui")],
}