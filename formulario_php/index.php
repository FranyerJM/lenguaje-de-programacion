<!DOCTYPE HTML> 
<html>
<head>
<style>
.error {color: #FF0000;}
.success {color: #008000; font-weight: bold;} 
</style>
</head>
<body> 

<h2>Formulario de Registro (by Franyer)</h2>
<p><span class="error">* campo obligatorio</span></p>

<form method="post" action="procesar.php"> 
  
  Nombre: <input type="text" name="nombre">
  <span class="error">*</span>
  <br><br>

  E-mail: <input type="text" name="email">
  <span class="error">*</span>
  <br><br>

  Contraseña: <input type="password" name="password">
  <span class="error">*</span>
  <br><br>

  <input type="submit" name="Registrarse" value="Registrarse"> 
</form>

</body>
</html>