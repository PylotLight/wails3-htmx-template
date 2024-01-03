/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["../components/*.go", "./dist/index.html", './src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}',],
    theme: {
        extend: {},
    },
    plugins: [require("daisyui")],
};