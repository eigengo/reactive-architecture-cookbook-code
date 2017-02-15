package com.reactivearchitecturecookbook.push

import akka.http.scaladsl.model.HttpRequest
import akka.http.scaladsl.server._
import akka.http.scaladsl.unmarshalling.GenericUnmarshallers
import com.reactivearchitecturecookbook.ingest.v1m0.IngestedImage

import scala.concurrent.{ExecutionContext, Future}

trait IngestRoute extends Directives with GenericUnmarshallers {

  def authorizeJwt(token: String): Boolean

  def extractIngestImage(request: HttpRequest)(implicit ec: ExecutionContext): Future[IngestedImage]

  def publishIngestedImage(token: String, transactionId: String)(ingestedImage: IngestedImage)(implicit ec: ExecutionContext): Future[Unit]

  private def authorizedRouteHandler(requestContext: RequestContext, token: String, transactionId: String)(implicit ec: ExecutionContext): Future[RouteResult] = {
    extractIngestImage(requestContext.request)
      .flatMap(publishIngestedImage(token, transactionId))
      .flatMap { _ ⇒ requestContext.complete("") }
      .recoverWith { case x ⇒ requestContext.fail(x) }
  }

  def ingestRoute(implicit ec: ExecutionContext): Route =
    path(RemainingPath) { transactionId ⇒
      post {
        extractCredentials {
          case Some(httpCredentials) if httpCredentials.scheme() == "Bearer" ⇒
            request: RequestContext ⇒
            if (authorizeJwt(httpCredentials.token())) authorizedRouteHandler(request, httpCredentials.token(), transactionId.toString())
            else request.reject(AuthorizationFailedRejection)
          case _ ⇒ _.reject(AuthorizationFailedRejection)
        }
      }
    }

}
