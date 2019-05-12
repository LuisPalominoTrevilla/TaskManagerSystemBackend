import com.tms.app._
import org.scalatra._
import javax.servlet.ServletContext
import scalaj.http.Http

class ScalatraBootstrap extends LifeCycle {
  override def init(context: ServletContext) {
    this.registerService
    context.mount(new AccountsController, "/*")
  }

  def registerService() {
    Http("http://127.0.0.1:3000/register").postForm(Seq("port" -> "4000", "service" -> "accounts", "healthCheck" -> "/healthCheck")).asString
  }
}
