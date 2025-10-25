mod tasks;

use chrono::{Datelike, Local, NaiveDate, TimeZone};
use clap::{Parser, Subcommand};
use std::{env, io};
use tasks::Tasks;

#[derive(Parser)]
#[clap(
    version = "1.0",
    author = "Franyer Marin",
    about = "todo simple"
)]
struct Cli {
    #[clap(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    Agregar {
        descripcion: Vec<String>,
    },
    Completar {
        numero: usize,
    },

    Progreso {
        numero: usize,
    },

    Modificar {

        numero: usize,
    },

    Eliminar {

        numero: usize,
    },

    Listar,
}


fn main() {
    let cli = Cli::parse();
    let mut tasks = Tasks::new();
    let filename = env::var("TODO_FILENAME")
        .unwrap_or_else(|_| ".todos.json".to_string());

    if let Err(e) = tasks.load(&filename) { // Pasa &filename
            eprintln!("error al cargar las tareas: {}", e);
            return;
        }

    match cli.command {
        Commands::Agregar { descripcion } => {
            let task_description = if descripcion.is_empty() {
                println!("introduce la descripcion de la tarea:");
                let mut buffer = String::new();
                io::stdin().read_line(&mut buffer).unwrap();
                buffer.trim().to_string()
            } else {
                descripcion.join(" ")
            };

            println!("introduce la fecha de entrega (dd/mm) o presiona Enter para no poner fecha:");
            let mut due_date_str = String::new();
            io::stdin().read_line(&mut due_date_str).unwrap();
            
            let due_date = if due_date_str.trim().is_empty() {
                None
            } else {
                let parts: Vec<&str> = due_date_str.trim().split('/').collect();
                if parts.len() == 2 {
                    if let (Ok(day), Ok(month)) = (parts[0].parse::<u32>(), parts[1].parse::<u32>()) {
                        let now = Local::now();
                        let year = now.year();
                        match NaiveDate::from_ymd_opt(year, month, day) {
                            Some(date) => Some(
                                Local
                                    .from_local_datetime(&date.and_hms_opt(23, 59, 59).unwrap())
                                    .unwrap(),
                            ),
                            None => {
                                eprintln!("fecha invalida (e.g., día 30 para febrero).");
                                return;
                            }
                        }
                    } else {
                        eprintln!("formato de fecha invalido. Asegurate de que sean numeros.");
                        return;
                    }
                } else {
                    eprintln!("formato de fecha invalido. Usa dd/mm.");
                    return;
                }
            };
            tasks.add(task_description, due_date);
        }
        Commands::Completar { numero } => {
            if let Err(e) = tasks.complete(numero) {
                eprintln!("{}", e);
            }
        }
        Commands::Progreso { numero } => {
            if let Err(e) = tasks.progress(numero) {
                eprintln!("{}", e);
            }
        }
        Commands::Modificar { numero } => {
            println!("introduce el nuevo titulo (deja en blanco para no cambiar):");
            let mut new_title = String::new();
            io::stdin().read_line(&mut new_title).unwrap();

            println!("introduce la nueva fecha de entrega (dd/mm), '-' para eliminar, o Enter para no cambiar:");
            let mut new_date_str = String::new();
            io::stdin().read_line(&mut new_date_str).unwrap();

            let new_due_date = match new_date_str.trim() {
                "" => None,
                "-" => Some(None),
                s => {
                    let parts: Vec<&str> = s.split('/').collect();
                    if parts.len() == 2 {
                        if let (Ok(day), Ok(month)) = (parts[0].parse::<u32>(), parts[1].parse::<u32>()) {
                            let now = Local::now();
                            let year = now.year();
                            match NaiveDate::from_ymd_opt(year, month, day) {
                                Some(date) => Some(Some(
                                    Local
                                        .from_local_datetime(&date.and_hms_opt(23, 59, 59).unwrap())
                                        .unwrap(),
                                )),
                                None => {
                                    eprintln!("fecha invalida (e.g., día 30 para febrero).");
                                    return;
                                }
                            }
                        } else {
                            eprintln!("formato de fecha invalido. Asegurate de que sean numeros.");
                            return;
                        }
                    } else {
                        eprintln!("formato de fecha invalido. Usa dd/mm.");
                        return;
                    }
                }
            };

            if let Err(e) = tasks.modify(
                numero,
                if new_title.trim().is_empty() {
                    None
                } else {
                    Some(new_title.trim().to_string())
                },
                new_due_date,
            ) {
                eprintln!("{}", e);
            }
        }
        Commands::Eliminar { numero } => {
            if let Err(e) = tasks.delete(numero) {
                eprintln!("{}", e);
            }
        }
        Commands::Listar => {
            tasks.print();
        }
    }

    if let Err(e) = tasks.store(&filename) { // Pasa &filename
            eprintln!("Error al guardar las tareas: {}", e);
        }
}