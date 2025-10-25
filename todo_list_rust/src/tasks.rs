use chrono::{Local, DateTime};
use colored::*;
use prettytable::{Cell, Row, Table};
use serde::{Deserialize, Serialize};
use std::fs::File;
use std::io::{Read, Write};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub enum TaskStatus {
    Pendiente,
    Progreso,
    Completada,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Task {
    #[serde(rename = "Tarea")]
    pub description: String,
    #[serde(rename = "Estado")]
    pub status: TaskStatus,
    #[serde(rename = "Creado_en")]
    pub created_at: DateTime<Local>,
    #[serde(rename = "Fecha_entrega")]
    #[serde(with = "custom_date_format")]
    pub due_date: Option<DateTime<Local>>,
    #[serde(rename = "Completado_en")]
    #[serde(with = "custom_date_format")]
    pub completed_at: Option<DateTime<Local>>,
}

mod custom_date_format {
    use chrono::{DateTime, Local};
    use serde::{self, Deserialize, Deserializer, Serializer};

    const FORMAT: &'static str = "%Y-%m-%dT%H:%M:%S.%f%z";
    const ALT_FORMAT: &'static str = "%Y-%m-%dT%H:%M:%S%z";

    pub fn serialize<S>(
        date: &Option<DateTime<Local>>,
        serializer: S,
    ) -> Result<S::Ok, S::Error>
    where
        S: Serializer,
    {
        if let Some(date) = date {
            let s = format!("{}", date.format(FORMAT));
            serializer.serialize_str(&s)
        } else {
            serializer.serialize_none()
        }
    }

    pub fn deserialize<'de, D>(
        deserializer: D,
    ) -> Result<Option<DateTime<Local>>, D::Error>
    where
        D: Deserializer<'de>,
    {
        let s: Option<String> = Option::deserialize(deserializer)?;
        if let Some(s) = s {
            if s == "0001-01-01T00:00:00Z" {
                return Ok(None);
            }
            if let Ok(dt) = DateTime::parse_from_str(&s, FORMAT) {
                return Ok(Some(dt.with_timezone(&Local)));
            }
            if let Ok(dt) = DateTime::parse_from_str(&s, ALT_FORMAT) {
                return Ok(Some(dt.with_timezone(&Local)));
            }
            Err(serde::de::Error::custom(format!("Error parsing date: {}", s)))
        } else {
            Ok(None)
        }
    }
}


pub struct Tasks(pub Vec<Task>);

impl Tasks {
    pub fn new() -> Self {
        Tasks(Vec::new())
    }

    pub fn add(&mut self, description: String, due_date: Option<DateTime<Local>>) {
        let task = Task {
            description,
            status: TaskStatus::Pendiente,
            created_at: Local::now(),
            due_date,
            completed_at: None,
        };
        self.0.push(task);
    }

    pub fn complete(&mut self, index: usize) -> Result<(), String> {
        if index == 0 || index > self.0.len() {
            return Err("indice incorrecto".to_string());
        }
        self.0[index - 1].status = TaskStatus::Completada;
        self.0[index - 1].completed_at = Some(Local::now());
        Ok(())
    }

    pub fn progress(&mut self, index: usize) -> Result<(), String> {
        if index == 0 || index > self.0.len() {
            return Err("indice incorrecto".to_string());
        }
        self.0[index - 1].status = TaskStatus::Progreso;
        Ok(())
    }

    pub fn modify(
        &mut self,
        index: usize,
        new_description: Option<String>,
        new_due_date: Option<Option<DateTime<Local>>>,
    ) -> Result<(), String> {
        if index == 0 || index > self.0.len() {
            return Err("indice incorrecto".to_string());
        }
        if let Some(desc) = new_description {
            if !desc.trim().is_empty() {
                self.0[index - 1].description = desc;
            }
        }
        if let Some(date) = new_due_date {
            self.0[index - 1].due_date = date;
        }
        Ok(())
    }

    pub fn delete(&mut self, index: usize) -> Result<(), String> {
        if index == 0 || index > self.0.len() {
            return Err("indice incorrecto".to_string());
        }
        self.0.remove(index - 1);
        Ok(())
    }

    pub fn load(&mut self, filename: &str) -> Result<(), String> {
        if let Ok(mut file) = File::open(filename) {
            let mut contents = String::new();
            file.read_to_string(&mut contents)
                .map_err(|e| e.to_string())?;
            if !contents.is_empty() {
                self.0 = serde_json::from_str(&contents).map_err(|e| e.to_string())?;
            }
        }
        Ok(())
    }

    pub fn store(&self, filename: &str) -> Result<(), String> {
        let mut file = File::create(filename).map_err(|e| e.to_string())?;
        let json = serde_json::to_string_pretty(&self.0).map_err(|e| e.to_string())?;
        file.write_all(json.as_bytes())
            .map_err(|e| e.to_string())?;
        Ok(())
    }

    pub fn print(&self) {
        let mut table = Table::new();
        table.add_row(Row::new(vec![
            Cell::new("#").style_spec("c"),
            Cell::new("Tarea").style_spec("c"),
            Cell::new("listo?").style_spec("c"),
            Cell::new("Fecha Creación").style_spec("r"),
            Cell::new("Fecha Entrega").style_spec("r"),
            Cell::new("Fecha Completado").style_spec("r"),
        ]));

        for (i, task) in self.0.iter().enumerate() {
            let (task_desc, status_str) = match task.status {
                TaskStatus::Completada => (
                    task.description.green(),
                    "si".green(),
                ),
                TaskStatus::Progreso => (
                    task.description.yellow(),
                    "empezó".yellow(),
                ),
                TaskStatus::Pendiente => {
                    let desc = if task.due_date.is_some() && task.due_date.unwrap() < Local::now() {
                        task.description.red()
                    } else {
                        task.description.blue()
                    };
                    (desc, "no".blue())
                }
            };

            let due_date_str = task.due_date.map_or("-".bright_black().to_string(), |d| d.format("%d/%m/%Y").to_string());
            let completed_at_str = task.completed_at.map_or("-".bright_black().to_string(), |d| d.format("%d/%m/%Y").to_string());

            table.add_row(Row::new(vec![
                Cell::new(&(i + 1).to_string()),
                Cell::new(&task_desc.to_string()),
                Cell::new(&status_str.to_string()),
                Cell::new(&task.created_at.format("%d/%m/%Y").to_string()),
                Cell::new(&due_date_str),
                Cell::new(&completed_at_str),
            ]));
        }

        table.printstd();
        println!(
            "{}",
            format!(
                "Tareas pendientes: {}                         Tareas atrasadas: {}",
                self.count_pending(),
                self.count_overdue()
            )
            .red()
        );
    }

    fn count_pending(&self) -> usize {
        self.0
            .iter()
            .filter(|t| t.status == TaskStatus::Pendiente)
            .count()
    }

    fn count_overdue(&self) -> usize {
        self.0
            .iter()
            .filter(|t| {
                t.status != TaskStatus::Completada
                    && t.due_date.is_some()
                    && t.due_date.unwrap() < Local::now()
            })
            .count()
    }
}

// TESTING 

#[cfg(test)]
mod tests {
    use super::*;
    use chrono::{Duration, Local};
    use std::fs;

    fn setup_tasks() -> Tasks {
        let tasks = Tasks(vec![
            Task {
                description: "Tarea 1".to_string(),
                status: TaskStatus::Pendiente,
                created_at: Local::now(),
                due_date: Some(Local::now() + Duration::days(1)),
                completed_at: None,
            },
            Task {
                description: "Tarea 2 (atrasada)".to_string(),
                status: TaskStatus::Pendiente,
                created_at: Local::now(),
                due_date: Some(Local::now() - Duration::days(1)),
                completed_at: None,
            },
            Task {
                description: "Tarea 3 (progreso)".to_string(),
                status: TaskStatus::Progreso,
                created_at: Local::now(),
                due_date: Some(Local::now() + Duration::days(2)),
                completed_at: None,
            },
        ]);
        tasks
    }

    #[test]
    fn test_add() {
        let mut tasks = Tasks::new();
        let due_date = Some(Local::now() + Duration::days(1));
        tasks.add("nueva tarea".to_string(), due_date);

        assert_eq!(tasks.0.len(), 1);
        assert_eq!(tasks.0[0].description, "nueva tarea");
        assert_eq!(tasks.0[0].status, TaskStatus::Pendiente);
    }

    #[test]
    fn test_complete() {
        let mut tasks = setup_tasks();

        // Si funciona
        let result = tasks.complete(1);
        assert!(result.is_ok());
        assert_eq!(tasks.0[0].status, TaskStatus::Completada);
        assert!(tasks.0[0].completed_at.is_some());

        // Si no funciona (indice invalido)
        let err_result = tasks.complete(99);
        assert!(err_result.is_err());
        assert_eq!(err_result.unwrap_err(), "indice incorrecto");
    }


    #[test]
    fn test_progress() {
        let mut tasks = setup_tasks();

        // Si funciona
        let result = tasks.progress(1);
        assert!(result.is_ok());
        assert_eq!(tasks.0[0].status, TaskStatus::Progreso);

        // Si no funciona (indice 0)
        let err_result = tasks.progress(0);
        assert!(err_result.is_err());
    }


    #[test]
    fn test_modify() {
        let mut tasks = setup_tasks();
        let new_date = Some(Some(Local::now() + Duration::days(3)));

        tasks.modify(1, Some("titulo modificado".to_string()), None).unwrap();
        assert_eq!(tasks.0[0].description, "titulo modificado");

        tasks.modify(1, None, new_date).unwrap();
        assert_eq!(tasks.0[0].due_date, new_date.unwrap());
        
        tasks.modify(2, Some("ambos modificados".to_string()), new_date).unwrap();
        assert_eq!(tasks.0[1].description, "ambos modificados");
        assert_eq!(tasks.0[1].due_date, new_date.unwrap());

        tasks.modify(1, None, Some(None)).unwrap();
        assert!(tasks.0[0].due_date.is_none());

        // indice invalido
        let err_result = tasks.modify(100, Some("no existe".to_string()), None);
        assert!(err_result.is_err());
    }

    #[test]
    fn test_delete() {
        let mut tasks = setup_tasks();
        let original_len = tasks.0.len();

        // Si funciona
        tasks.delete(1).unwrap();
        assert_eq!(tasks.0.len(), original_len - 1);
        assert_eq!(tasks.0[0].description, "Tarea 2 (atrasada)");

        // Si no funciona
        assert!(tasks.delete(99).is_err());
    }

    // Equivalente a TestAlmacenarYCargar
    #[test]
    fn test_store_and_load() {
        let tasks = setup_tasks();
        let filename = ".test_tasks_storage.json";

        // Almacenar
        let store_result = tasks.store(filename);
        assert!(store_result.is_ok());

        // Cargar
        let mut loaded_tasks = Tasks::new();
        let load_result = loaded_tasks.load(filename);
        assert!(load_result.is_ok());

        // Verificar
        assert_eq!(tasks.0.len(), loaded_tasks.0.len());
        assert_eq!(tasks.0[0].description, loaded_tasks.0[0].description);
        assert_eq!(tasks.0[1].status, loaded_tasks.0[1].status);

        // Limpieza
        fs::remove_file(filename).unwrap();
    }

    #[test]
    fn test_count_pending() {
        let tasks = setup_tasks();
        // Tarea 1 y Tarea 2 son Pendientes
        assert_eq!(tasks.count_pending(), 2);
    }


    #[test]
    fn test_count_overdue() {
        let tasks = setup_tasks();
        // Tarea 2 esta atrasada
        assert_eq!(tasks.count_overdue(), 1);
        
        // Si la completamos, ya no debe contar
        let mut tasks_completada = setup_tasks();
        tasks_completada.complete(2).unwrap();
        assert_eq!(tasks_completada.count_overdue(), 0);
    }
}