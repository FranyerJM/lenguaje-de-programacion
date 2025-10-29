<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Mi Gatito 🐾</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Pacifico&family=Quicksand:wght@300;400;700&display=swap" rel="stylesheet">
    <style>
        body {
            font-family: 'Quicksand', sans-serif;
            margin: 0;
            background-color: #fce4ec
            color: #4a4a4a;
            line-height: 1.6;
        }

        h1, h2 {
            font-family: 'Pacifico', ;
            color: #e91e63;
            text-align: center;
            margin-bottom: 20px;
        }

        header {
            background-color: #ff80ab;
            color: white;
            padding: 20px 0;
            text-align: center;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }

        header img {
            border-radius: 50%;
            width: 150px;
            height: 150px;
            object-fit: cover;
            border: 5px solid #f8bbd0;
            margin-bottom: 10px;
        }

        header h1 {
            font-size: 3em;
            margin: 0;
            color: white;
        }

        .about-me {
            max-width: 800px;
            margin: 30px auto;
            padding: 20px;
            background-color: white;
            border-radius: 15px;
            box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
            text-align: center;
        }

        .about-me p {
            font-size: 1.1em;
            color: #6a6a6a;
        }

        .gallery-container {
            max-width: 1000px;
            margin: 40px auto;
            padding: 20px;
            background-color: white;
            border-radius: 15px;
            box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
        }

        .gallery {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 15px;
        }

        .gallery img {
            width: 100%;
            height: 250px;
            object-fit: cover
            border-radius: 10px;
            border: 3px solid #ffab91;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
            transition: transform 0.3s ease-in-out;
        }

        .gallery img:hover {
            transform: scale(1.05);
        }


        footer {
            text-align: center;
            padding: 20px;
            margin-top: 40px;
            background-color: #ff80ab;
            color: white;
            font-size: 0.9em;
        }


        @media (max-width: 768px) {
            header h1 {
                font-size: 2.5em;
            }
            .about-me, .gallery-container {
                margin: 20px 10px;
                padding: 15px;
            }
        }
    </style>
</head>
<body>

    <header>
        <img src="img/image1.JPG" alt="Foto principal de mi gatito">
        <h1>Fotos de mi Gatito Lindo 🥹❤️!</h1> 
    </header>

    <main>
        <section class="about-me">
            <h2>Sobre Mi</h2>
            <p>Soy un gatito muy lindo tierno, adorable, hermoso, inteligente, y tengo una casa muy grande donde me dan mucho amor todos los dias, me hacen sentir en la gatito mas feliz del gatimundo</p>
            <p>tengo 1y de edad y desde que naci, lleno de ternura a todo el humano que me ve, tan solo mira mis fotos 🥰</p>
        </section>

        <section class="gallery-container">
            <h2>Gati Moments</h2>
            <div class="gallery">
                <img src="img/image1.JPG">
                <img src="img/image2.JPG">
                <img src="img/image3.JPG">
                <img src="img/image4.JPG">
                <img src="img/image5.HEIC">
                <img src="img/image6.HEIC">
                <img src="img/image7.HEIC">
                <img src="img/image8.HEIC">
                <img src="img/image9.HEIC">
            </div>
        </section>
    </main>

    <footer>
        <p>Hecho con mucho amor para mi gatito💖 | © <?php echo date("Y"); ?> Fran</p>
    </footer>

</body>
</html>