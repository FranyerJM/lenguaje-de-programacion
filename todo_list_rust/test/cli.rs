use assert_cmd::Command;
use chrono::{DateTime, Duration, Local};
use predicates::prelude::*;
use serde_json::json;
use std::env;
use std::fs;
use std::io::Write;
use std::path::PathBuf;


fn get_test_filename() -> PathBuf {
    let mut path = env::temp_dir();
    path.push(".test_todos.json");
    path
}


fn setup_test_file(path: &PathBuf) {

    let test_task = json!([{
        "Tarea": "Tarea de prueba 1",
        "Estado": "Pendiente",
        "Creado_en": Local::now(),
        "Fecha_entrega": Local::now() + Duration::days(1),
        "Completado_en": null
    }]);

    let mut file = fs::File::create(path).expect("No se pudo crear el archivo de prueba");
    file.write_all(test_task.to_string().as_bytes())
        .expect("No se pudo escribir al archivo de prueba");
}

// Limpieza
fn cleanup_file(path: &PathBuf) {
    let _ = fs::remove_file(path);
}


fn get_cmd() -> Command {
    Command::cargo_bin(env!("CARGO_PKG_NAME")).unwrap()
}


#[test]
fn test_agregar_nueva_tarea() {
    let filename = get_test_filename();
    cleanup_file(&filename);
    env::set_var("TODO_FILENAME", filename.as_os_str());

    let mut cmd = get_cmd();
    cmd.arg("agregar")
        .arg("nueva tarea desde prueba")
        .write_stdin("\n") // Simula Enter para la fecha
        .assert()
        .success()
        .stdout(predicate::str::contains("introduce la fecha de entrega"));

    cleanup_file(&filename);
}


#[test]
fn test_listar_tareas() {
    let filename = get_test_filename();
    setup_test_file(&filename);
    env::set_var("TODO_FILENAME", filename.as_os_str());

    let mut cmd = get_cmd();
    cmd.arg("listar")
        .assert()
        .success()
        .stdout(predicate::str::contains("Tarea de prueba 1"));

    cleanup_file(&filename);
}


#[test]
fn test_completar_tarea() {
    let filename = get_test_filename();
    setup_test_file(&filename);
    env::set_var("TODO_FILENAME", filename.as_os_str());

    let mut cmd_complete = get_cmd();
    cmd_complete.arg("completar").arg("1").assert().success();

    let mut cmd_list = get_cmd();
    cmd_list.arg("listar")
        .assert()
        .success()
        .stdout(predicate::str::contains("\x1b[32msi\x1b[0m")); 

    cleanup_file(&filename);
}

#[test]
fn test_eliminar_tarea() {
    let filename = get_test_filename();
    setup_test_file(&filename);
    env::set_var("TODO_FILENAME", filename.as_os_str());

    let mut cmd_delete = get_cmd();
    cmd_delete.arg("eliminar").arg("1").assert().success();

    let content = fs::read_to_string(&filename).unwrap_or_default();
    assert_eq!(content, "[]", "El archivo no quedó vacío después de eliminar");

    cleanup_file(&filename);
}

#[test]
fn test_sin_argumentos_muestra_ayuda() {
    let filename = get_test_filename();
    env::set_var("TODO_FILENAME", filename.as_os_str());

    let mut cmd = get_cmd();

    cmd.assert()
        .failure()

        .stderr(predicate::str::contains("todo simple")); 
}

#[test]
fn test_modificar_tarea() {
    let filename = get_test_filename();
    setup_test_file(&filename);
    env::set_var("TODO_FILENAME", filename.as_os_str());

    let mut cmd = get_cmd();
    cmd.arg("modificar")
        .arg("1")
        .write_stdin("Tarea modificada\n-\n")
        .assert()
        .success();

    // Verificar
    let mut cmd_list = get_cmd();
    cmd_list.arg("listar")
        .assert()
        .success()
        .stdout(predicate::str::contains("Tarea modificada"));

    cleanup_file(&filename);
}