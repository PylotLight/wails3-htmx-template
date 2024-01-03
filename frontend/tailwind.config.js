/** @type {import('tailwindcss').Config} */
export default {
    content: ["../components/*.go", "./index.html", './src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}',],
    theme: {
        extend: {},
    },
    plugins: [require("daisyui")],
};