<?php

$nombreError = $emailError = $passwordError = "";
$nombre = $email = "";

function limpiar($data) {
  $data = trim($data);
  $data = stripslashes($data);
  $data = htmlspecialchars($data);
  return $data;
}

if ($_SERVER["REQUEST_METHOD"] == "POST") {

  
  if (empty($_POST["nombre"])) {
    $nombreError = "el Nombre es obligatorio";
  } else {
    $nombre = limpiar($_POST["nombre"]);
    if (!preg_match("/^[a-zA-Z-' ]*$/", $nombre)) {
      $nombreError = "Solo se permiten letras y espacios en blanco";
    }
  }
  
  if (empty($_POST["email"])) {
    $emailError = "el Email es obligatorio";
  } else {
    $email = limpiar($_POST["email"]);
    if (!filter_var($email, FILTER_VALIDATE_EMAIL)) {
      $emailError = "el email no tiene un formato correcto";
    }
  }
    
  if (empty($_POST["password"])) {
    $passwordError = "La Contraseña es obligatoria";
  }

  
  if ($nombreError == "" && $emailError == "" && $passwordError == "") {
      
      $passOriginal = trim($_POST["password"]);
      $passHash = password_hash($passOriginal, PASSWORD_DEFAULT);

      $nuevoUsuario = [
          "nombre" => $nombre,
          "email" => $email,
          "password" => $passHash
      ];

      $archivo = 'usuarios.json';
      
      $usuarios = [];
      if (file_exists($archivo)) {
          $jsonExistente = file_get_contents($archivo);
          $usuarios = json_decode($jsonExistente, true);
      }

      $usuarios[] = $nuevoUsuario;
      $json_final = json_encode($usuarios, JSON_PRETTY_PRINT);
      
      if (file_put_contents($archivo, $json_final)) {

          echo "<h1>Usuario registrado correctamente!</h1>";
          echo '<p><a href="index.php">Volver al formulario</a></p>';
      } else {
          // ERROR DE ESCRITURA
          echo "<h1>Error</h1>";
          echo "<p>No se pudo guardar en el archivo. Verifica los permisos.</p>";
          echo '<p><a href="index.php">Volver a intentarlo</a></p>';
      }

  } else {

      echo "<h1>Error en el formulario</h1>";
      echo "<p>Por favor, corrige lo siguiente:</p>";
      
      echo '<div style="color: #FF0000;">';
      if ($nombreError != "") echo $nombreError . "<br>";
      if ($emailError != "") echo $emailError . "<br>";
      if ($passwordError != "") echo $passwordError . "<br>";
      echo '</div>';
      
      echo '<p><a href="index.php">Volver a intentarlo</a></p>';
  }

} else {
    //seguridad: si se accede directamente a procesar.php sin pasar por el formulario
  echo "Acceso no permitido. Por favor, rellena el formulario.";
  echo '<p><a href="index.php">Ir al formulario</a></p>';
}

?>