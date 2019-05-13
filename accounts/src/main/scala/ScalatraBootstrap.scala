import com.tms.app._
import org.scalatra._
import javax.servlet.ServletContext
import scalaj.http.Http

class ScalatraBootstrap extends LifeCycle {
  override def init(context: ServletContext) {
    this.registerService
    context.mount(new AccountsController, "/*")
    context.initParameters("org.scalatra.cors.allowedOrigins") = "*"
    context.initParameters("org.scalatra.cors.allowCredentials") = "false"
  }

  def registerService() {
    val host = sys.env.get("REGISTRY_HOST").getOrElse("http://127.0.0.1")
    val endpoint = sys.env.get("REGISTRY_ENDPOINT").getOrElse("/register")

    try {
      val request = Http(host + endpoint).postForm(Seq("port" -> "4000", "service" -> "accounts", "healthCheck" -> "/healthCheck")).asString
      println(request.body);
    } catch {
      case e: Throwable => {
        println("Registry is not available")
      }
    }
  }
}
