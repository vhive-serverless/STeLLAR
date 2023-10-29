use actix_web::{get, web, App, HttpRequest, HttpServer, Responder, Result};
use serde::{Deserialize, Serialize};
use std::time::SystemTime;

#[derive(Debug, Deserialize)]
pub struct Params {
    #[serde(rename = "IncrementLimit")]
    #[serde(default)]
    increment_limit: u64,
}

#[derive(Serialize)]
struct Response {
    #[serde(rename = "RequestID")]
    request_id: String,
    #[serde(rename = "TimestampChain")]
    timestamp_chain: Vec<String>,
}

fn simulate_work(increment_limit: u64) {
    let mut i = 0;
    while i < increment_limit {
        i += 1;
    }
}

fn get_system_time() -> String {
    let mut buffer = itoa::Buffer::new();
    match SystemTime::now().duration_since(SystemTime::UNIX_EPOCH) {
        Ok(n) => buffer.format(n.as_secs()).to_owned(),
        Err(_) => panic!("SystemTime before UNIX EPOCH!"),
    }
}

#[get("/")]
async fn hellorust(req: HttpRequest) -> Result<impl Responder> {
    let params = web::Query::<Params>::from_query(req.query_string()).unwrap();
    if params.increment_limit > 0 {
        simulate_work(params.increment_limit)
    }

    let res_obj = Response {
        request_id: "google-does-not-specify".to_owned(),
        timestamp_chain: vec![get_system_time()],
    };
    Ok(web::Json(res_obj))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| App::new().service(hellorust))
        .bind(("0.0.0.0", 8080))?
        .run()
        .await
}
