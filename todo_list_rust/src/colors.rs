use colored::*;

pub fn red(s: &str) -> String {
    s.red().to_string()
}

pub fn green(s: &str) -> String {
    s.green().to_string()
}

pub fn blue(s: &str) -> String {
    s.blue().to_string()
}

pub fn grey(s: &str) -> String {
    s.truecolor(150, 150, 150).to_string()
}

pub fn yellow(s: &str) -> String {
    s.yellow().to_string()
}

// TESTING

// Añade esto al final de src/colors.rs

#[cfg(test)]
mod tests {
    use super::*;
    use colored::control;

    // Helper para comparar la salida ANSI, similar a tu test de Go
    fn assert_color(
        color_fn: fn(&str) -> String,
        text: &str,
        expected_ansi: &str,
    ) {
        // Forzamos que la librería 'colored' emita códigos ANSI para el test
        control::set_override(true);
        let result = color_fn(text);
        control::unset_override(); // Lo reseteamos
        
        // Comparamos que la salida contenga el código ANSI esperado
        assert!(
            result.contains(expected_ansi),
            "Esperado: ...{}..., Obtenido: {}",
            expected_ansi,
            result
        );
        assert!(
            result.contains(text),
            "La salida no contiene el texto original"
        );
    }

    // Equivalente a TestRojo
    #[test]
    fn test_red() {
        assert_color(red, "texto", "\x1b[31m");
    }

    // Equivalente a TestVerde
    #[test]
    fn test_green() {
        assert_color(green, "texto", "\x1b[32m");
    }

    // Equivalente a TestAzul
    #[test]
    fn test_blue() {
        assert_color(blue, "texto", "\x1b[34m");
    }
    
    // Equivalente a TestAmarillo
    #[test]
    fn test_yellow() {
        assert_color(yellow, "texto", "\x1b[33m");
    }

    // Test para 'grey' (el único con lógica custom)
    // Equivalente a TestGris (aunque tu Go usa un 'gris' simple)
    #[test]
    fn test_grey() {
        // Este usa truecolor
        assert_color(grey, "texto", "\x1b[38;2;150;150;150m");
    }
}